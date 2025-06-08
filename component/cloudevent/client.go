package cloudevent

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/nats-io/nats.go/jetstream"
	ants "github.com/panjf2000/ants/v2"
	"github.com/pubgo/funk/anyhow"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/component/natsclient"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/errors/errcheck"
	"github.com/pubgo/funk/log"
	cloudeventpb "github.com/pubgo/funk/proto/cloudevent"
	"github.com/pubgo/funk/running"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
	"github.com/pubgo/funk/typex"
	"github.com/pubgo/funk/version"
	"github.com/pubgo/lava/core/lifecycle"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Params struct {
	Nc  *natsclient.Client
	Cfg *Config
	Lc  lifecycle.Lifecycle
}

func New(p Params) *Client {
	js := assert.Must1(jetstream.New(p.Nc.Conn))
	return &Client{
		p:         p,
		js:        js,
		prefix:    DefaultPrefix,
		handlers:  make(map[string]map[string]EventHandler[proto.Message]),
		streams:   make(map[string]jetstream.Stream),
		consumers: make(map[string]map[string]*Consumer),
		jobs:      make(map[string]map[string]map[string]*jobEventHandler),
		subjects:  getAllSubject(),
	}
}

type Client struct {
	p  Params
	js jetstream.JetStream

	// stream manager
	streams map[string]jetstream.Stream

	// jobs: stream->consumer->Consumer
	consumers map[string]map[string]*Consumer

	// handlers: job name -> subject -> job handler
	handlers map[string]map[string]EventHandler[proto.Message]

	// jobs: stream->consumer->subject->jobEventHandler
	jobs map[string]map[string]map[string]*jobEventHandler

	// stream, consumer, subject prefix, default: DefaultPrefix
	prefix string

	// subjects operation => subject info
	subjects map[string]*cloudeventpb.CloudEventMethodOptions
}

func (c *Client) initStream() (r error) {
	defer errcheck.RecoveryAndCheck(&r)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	for streamName, cfg := range c.p.Cfg.Streams {
		streamName = c.streamName(streamName)

		assert.If(c.streams[streamName] != nil, "stream %s already exists", streamName)

		// add subject prefix
		streamSubjects := lo.Map(cfg.Subjects, func(item string, index int) string { return c.subjectName(item) })
		metadata := map[string]string{"creator": fmt.Sprintf("%s/%s/%s", version.Project(), version.Version(), running.InstanceID)}
		storageType := getStorageType(cfg.Storage)
		streamCfg := jetstream.StreamConfig{
			Name:     streamName,
			Subjects: streamSubjects,
			Metadata: metadata,
			Storage:  storageType,
			//Retention: jetstream.InterestPolicy,

			// Duplicates is the window within which to track duplicate messages.
			// If not set, server default is 2 minutes.
			Duplicates: time.Minute * 5,
		}

		stream, err := c.js.CreateOrUpdateStream(ctx, streamCfg)
		err = errors.IfErr(err, func(err error) error {
			return errors.Wrapf(err, "failed to create stream:%s", streamName)
		})
		if errcheck.Check(&r, err) {
			return
		}
		c.streams[streamName] = stream
	}
	return
}

func (c *Client) initConsumer() (r error) {
	defer errcheck.RecoveryAndCheck(&r)

	allEventKeysSet := mapset.NewSet(lo.MapToSlice(c.subjects, func(key string, value *cloudeventpb.CloudEventMethodOptions) string { return c.subjectName(key) })...)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	for jobOrConsumerName, consumers := range c.p.Cfg.Consumers {
		jobName := jobOrConsumerName
		assert.If(c.handlers[jobName] == nil, "failed to find job handler: %s, please impl RegisterCloudJob", jobName)

		consumerName := jobOrConsumerName
		for _, cfg := range consumers {
			// check subject exists
			for _, sub := range cfg.Subjects {
				name := c.subjectName(lo.FromPtr(sub.Name))
				assert.If(!allEventKeysSet.Contains(name), "subject:%s not found, please check protobuf define and service", name)
			}

			consumerName = c.consumerName(lo.Ternary(cfg.Consumer != nil, lo.FromPtr(cfg.Consumer), consumerName))
			streamName := c.streamName(cfg.Stream)

			// consumer init
			typex.DoBlock(func() {
				if c.consumers[streamName] == nil {
					c.consumers[streamName] = make(map[string]*Consumer)
				}
				// A streaming consumer can only have one corresponding job handler
				assert.If(c.consumers[streamName][consumerName] != nil, "consumer %s already exists", consumerName)

				metadata := map[string]string{"version": fmt.Sprintf("%s/%s", version.Project(), version.Version())}
				consumerCfg := jetstream.ConsumerConfig{
					Name:     consumerName,
					Durable:  consumerName,
					Metadata: metadata,
				}

				consumer, err := c.js.CreateOrUpdateConsumer(ctx, streamName, consumerCfg)
				assert.Fn(err != nil, func() error {
					return errors.Wrapf(err, "stream=%s consumer=%s", streamName, consumerName)
				})
				logger.Info().Func(func(e *zerolog.Event) {
					e.Str("stream", streamName)
					e.Str("consumer", consumerName)
					e.Msg("register consumer success")
				})

				c.consumers[streamName][consumerName] = &Consumer{Consumer: consumer, Config: cfg}
			})

			typex.DoBlock(func() {
				if c.jobs[streamName] == nil {
					c.jobs[streamName] = make(map[string]map[string]*jobEventHandler)
				}

				if c.jobs[streamName][consumerName] == nil {
					c.jobs[streamName][consumerName] = map[string]*jobEventHandler{}
				}

				baseJobConfig := handleDefaultJobConfig(cfg.Job)
				subjectMap := lo.SliceToMap(cfg.Subjects, func(item1 *strOrJobConfig) (string, *JobEventConfig) {
					item := lo.ToPtr(JobEventConfig(lo.FromPtr(item1)))
					return c.subjectName(*item.Name), mergeJobConfig(item, baseJobConfig)
				})

				for subName, subCfg := range subjectMap {
					assert.If(c.handlers[jobName][subName] == nil, "job handler not found, job_name=%s sub_name=%s", jobName, subName)

					job := &jobEventHandler{
						name:    jobName,
						handler: c.handlers[jobName][subName],
						cfg:     subCfg,
					}

					logger.Info().Func(func(e *zerolog.Event) {
						e.Str("job_name", job.name)
						e.Str("job_handler", stack.CallerWithFunc(job.handler).String())
						e.Any("job_config", subCfg)
						e.Any("stream_name", streamName)
						e.Any("consumer_name", consumerName)
						e.Any("job_subject", subName)
						e.Msg("register cloud job handler executor")
					})
					c.jobs[streamName][consumerName][subName] = job
				}
			})
		}
	}
	return
}

func (c *Client) doConsumeHandler(streamName, consumerName string, jobSubjects map[string]*jobEventHandler, concurrent int) func(msg jetstream.Msg) {
	var handler = func(msg jetstream.Msg) {
		var now = time.Now()
		var addMsgInfo = func(e *zerolog.Event) {
			e.Str("stream", streamName)
			e.Str("consumer", consumerName)
			e.Any("header", msg.Headers())
			e.Any("msg_id", msg.Headers().Get(jetstream.MsgIDHeader))
			e.Str("subject", msg.Subject())
			e.Str("msg_received_time", now.String())
			e.Str("job_cost", time.Since(now).String())
		}

		logger.Debug().Func(func(e *zerolog.Event) {
			addMsgInfo(e)
			e.Msg("received cloud job event")
		})

		var handlerDelayJob = func() (_ bool, gErr error) {
			delayDur := strings.TrimSpace(msg.Headers().Get(DefaultCloudEventDelayKey))
			if delayDur == "" {
				return false, nil
			}

			dur := decodeDelayTime(delayDur).MapErr(func(err error) error {
				return errors.Wrap(err, "failed to parse cloud job delay time")
			})
			if dur.Catch(&gErr) {
				return
			}

			durVal := dur.GetValue()
			// ignore negative delay
			if durVal < 0 {
				return false, nil
			}

			return true, msg.NakWithDelay(durVal)
		}

		if ok, err := handlerDelayJob(); err != nil {
			logger.Err(err).Func(addMsgInfo).Msg("failed to handle cloud delay job and no ack")
			return
		} else if ok {
			logger.Info().Func(addMsgInfo).Msg("redeliver the message after the given delay")
			return
		}

		handler := jobSubjects[msg.Subject()]
		if handler == nil {
			logger.Error().Func(addMsgInfo).Msg("failed to find subject job handler")
			return
		}

		meta, err := msg.Metadata()
		if err != nil {
			// no ack, retry always, unless it can recognize special error information
			logger.Err(err).Func(addMsgInfo).Msg("failed to parse nats stream msg metadata")
			return
		}

		var cfg = handler.cfg
		var checkErrAndLog = func(err error, msg string) {
			if err == nil {
				return
			}

			logger.Err(err).
				Str("fn_caller", stack.Caller(1).String()).
				Func(addMsgInfo).
				Any("metadata", meta).
				Any("config", cfg).
				Any("msg_received_time", now.String()).
				Str("job_cost", time.Since(now).String()).
				Msg(msg)
		}

		err = try.Try(func() error { return c.doHandler(meta, msg, handler, cfg).GetErr() })
		if err == nil {
			checkErrAndLog(msg.Ack(), "failed to do msg ack with handler ok")
			return
		}

		// reject job msg
		if isRejectErr(err) {
			checkErrAndLog(msg.TermWithReason("reject by caller"), "failed to do msg ack with reject err")
			return
		}

		var backoff = lo.FromPtr(cfg.RetryBackoff)
		var maxRetries = lo.FromPtr(cfg.MaxRetry)

		// If the error is a redelivery error, then the backoff duration is the error duration
		if err1 := isRedeliveryErr(err); err1 != nil {
			backoff = err1.delay
		}

		// Proactively retry and did not reach the maximum retry count
		if meta.NumDelivered < uint64(maxRetries) {
			logger.Warn().
				Err(err).
				Func(addMsgInfo).
				Any("metadata", meta).
				Msg("retry nats stream cloud job event")
			checkErrAndLog(msg.NakWithDelay(backoff), "failed to retry msg with delay nak")
			return
		}

		checkErrAndLog(err, "failed to do handler cloud job")
		checkErrAndLog(msg.Ack(), "failed to do msg ack with handler error")
	}

	pool := assert.Must1(ants.NewPool(
		concurrent,
		ants.WithLogger(log.NewStd(logger)),
		ants.WithNonblocking(false),
	))
	// pool.Release()
	return func(msg jetstream.Msg) {
		if pool.Running() == concurrent {
			logger.Warn().Func(func(e *zerolog.Event) {
				e.Int("concurrent", concurrent)
				e.Str("stream", streamName)
				e.Str("consumer", consumerName)
				e.Msg("concurrent limit occurred, please check the concurrent limit")
			})
		}
		if err := pool.Submit(func() { handler(msg) }); err != nil {
			logger.Err(err).Func(func(e *zerolog.Event) {
				e.Str("stream", streamName)
				e.Str("consumer", consumerName)
				e.Msg("failed to submit job to pool")
			})
		}
	}
}

func (c *Client) doErrHandler(streamName, consumerName string) jetstream.PullConsumeOpt {
	return jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		logger.Err(err).
			Str("stream", streamName).
			Str("consumer", consumerName).
			Msg("nats consumer error")
	})
}

func (c *Client) doHandler(meta *jetstream.MsgMetadata, msg jetstream.Msg, job *jobEventHandler, cfg *JobEventConfig) (gErr anyhow.Error) {
	var timeout = lo.FromPtr(cfg.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx = log.UpdateEventCtx(ctx, log.Map{
		"sub_subject":                 msg.Subject(),
		"sub_stream":                  meta.Stream,
		"sub_consumer":                meta.Consumer,
		"sub_msg_id":                  msg.Headers().Get(jetstream.MsgIDHeader),
		"sub_msg_" + DefaultSenderKey: msg.Headers().Get(DefaultSenderKey),
	})

	msgCtx := &Context{
		Header:       http.Header(msg.Headers()),
		NumDelivered: meta.NumDelivered,
		NumPending:   meta.NumPending,
		Timestamp:    meta.Timestamp,
		Stream:       meta.Stream,
		Consumer:     meta.Consumer,
		Subject:      msg.Subject(),
		Config:       cfg,
	}

	var now = time.Now()
	var args any
	defer func() {
		if gErr.IsOK() {
			return
		}

		logger.Err(gErr.GetErr()).Func(func(e *zerolog.Event) {
			e.Any("context", msgCtx)
			e.Any("args", args)
			e.Str("timeout", timeout.String())
			e.Str("start_time", now.String())
			e.Str("job_cost", time.Since(now).String())
			e.Msg("failed to do cloud job handler")
		})
	}()

	var pb anypb.Any
	err := anyhow.ErrOf(proto.Unmarshal(msg.Data(), &pb)).
		MapErr(func(err error) error {
			return errors.WrapTag(err,
				errors.T("msg", "failed to unmarshal stream msg data to any proto"),
				errors.T("args", string(msg.Data())),
			)
		})
	if err.Catch(&gErr) {
		return
	}
	args = &pb

	dst := anyhow.Wrap(anypb.UnmarshalNew(args.(*anypb.Any), proto.UnmarshalOptions{})).
		WithErr(func(err error) error {
			return errors.WrapTag(err,
				errors.T("msg", "failed to unmarshal any proto to proto msg"),
				errors.T("args", args),
			)
		})
	if dst.Catch(&gErr) {
		return
	}

	ctx = createCtxWithContext(ctx, msgCtx)
	err = anyhow.ErrOf(job.handler(ctx, dst.GetValue())).
		MapErr(func(err error) error {
			return errors.WrapTag(err,
				errors.T("msg", "failed to do cloud job handler"),
				errors.T("args", dst),
				errors.T("any_pb", dst),
			)
		})
	return err
}

func (c *Client) doConsume() (r error) {
	defer errcheck.RecoveryAndCheck(&r)
	for streamName, consumers := range c.consumers {
		for consumerName, consumer := range consumers {
			assert.If(c.jobs[streamName] == nil, "stream not found, stream=%s", streamName)
			assert.If(c.jobs[streamName][consumerName] == nil, "consumer not found, consumer=%s", consumerName)

			jobSubjects := c.jobs[streamName][consumerName]

			concurrent := DefaultConcurrent
			if consumer.Config.Concurrent != nil {
				concurrent = lo.FromPtr(consumer.Config.Concurrent)
			}
			if concurrent < DefaultMinConcurrent || concurrent > DefaultMaxConcurrent {
				return errors.Errorf("concurrent must be in the range of %d-%d", DefaultMinConcurrent, DefaultMaxConcurrent)
			}

			logger.Info().Func(func(e *zerolog.Event) {
				e.Str("stream", streamName)
				e.Str("consumer", consumerName)
				e.Any("subjects", lo.MapKeys(jobSubjects, func(_ *jobEventHandler, key string) string { return key }))
				e.Msg("cloud job do consumer")
			})

			con := assert.Must1(consumer.Consume(
				c.doConsumeHandler(streamName, consumerName, jobSubjects, concurrent),
				c.doErrHandler(streamName, consumerName),
			))
			c.p.Lc.BeforeStop(func() { con.Stop() })
		}
	}
	return
}

func (c *Client) Start() error {
	assert.Exit(c.initStream())
	assert.Exit(c.initConsumer())
	assert.Exit(c.doConsume())
	return nil
}

func (c *Client) streamName(name string) string {
	prefix := fmt.Sprintf("%s:", c.prefix)
	if strings.HasPrefix(name, prefix) {
		return name
	}

	return fmt.Sprintf("%s%s", prefix, name)
}

func (c *Client) consumerName(name string) string {
	prefix := fmt.Sprintf("%s:", c.prefix)
	if strings.HasPrefix(name, prefix) {
		return name
	}

	return fmt.Sprintf("%s%s", prefix, name)
}

func (c *Client) subjectName(name string) string {
	return handleSubjectName(name, c.prefix)
}

func (c *Client) GetSubject(name string) *cloudeventpb.CloudEventMethodOptions {
	return c.subjects[name]
}
