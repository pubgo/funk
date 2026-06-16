package cloudevent

import (
	"context"
	"fmt"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/component/lifecycle"
	"github.com/pubgo/funk/v2/component/natsclient"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/running"
	"github.com/pubgo/funk/v2/stack"
	"github.com/pubgo/funk/v2/typex"
	cloudeventpb "github.com/pubgo/funk/v2/proto/cloudevent"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type Params struct {
	Nc  *natsclient.Client
	Cfg *Config
	Lc  lifecycle.Lifecycle
}

func New(p Params) *Client {
	js := assert.Must1(jetstream.New(p.Nc.Conn))
	return &Client{
		p:           p,
		js:          js,
		prefix:      DefaultPrefix,
		jobManagers: make(map[string]*jobManager),
		streams:     make(map[string]jetstream.Stream),
		consumers:   make(map[string]map[string]*Consumer),
		jobs:        make(map[string]map[string]map[string]*jobEventHandler),
		subjects:    getAllSubject(),
	}
}

type Client struct {
	p  Params
	js jetstream.JetStream

	streams   map[string]jetstream.Stream
	consumers map[string]map[string]*Consumer
	jobManagers map[string]*jobManager
	jobs      map[string]map[string]map[string]*jobEventHandler
	prefix    string
	subjects  map[string]*cloudeventpb.CloudEventMethodOptions
}

func (c *Client) initStream() (r result.Error) {
	defer result.Recovery(&r)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	for streamName, cfg := range c.p.Cfg.Streams {
		streamName = c.streamName(streamName)

		assert.If(c.streams[streamName] != nil, "stream %s already exists", streamName)

		streamSubjects := lo.Map(cfg.Subjects, func(item string, index int) string { return c.subjectName(item) })
		metadata := map[string]string{"creator": fmt.Sprintf("%s/%s/%s", version.Project(), version.Version(), running.InstanceID)}
		streamCfg := jetstream.StreamConfig{
			Name:       streamName,
			Subjects:   streamSubjects,
			Metadata:   metadata,
			Storage:    getStorageType(cfg.Storage),
			Duplicates: time.Minute * 5,
		}

		stream := result.Wrap(c.js.CreateOrUpdateStream(ctx, streamCfg)).
			MapErr(func(err error) error {
				return errors.Wrapf(err, "failed to create stream:%s", streamName)
			}).
			UnwrapOrThrow(&r)
		if r.IsErr() {
			return
		}
		c.streams[streamName] = stream
	}
	return
}

func (c *Client) initConsumer() (r result.Error) {
	defer result.Recovery(&r)

	allEventKeysSet := mapset.NewSet(lo.MapToSlice(c.subjects, func(key string, value *cloudeventpb.CloudEventMethodOptions) string { return c.subjectName(key) })...)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	for jobOrConsumerName, consumers := range c.p.Cfg.Consumers {
		jobName := jobOrConsumerName
		assert.If(c.jobManagers[jobName] == nil, "failed to find job manager: %s, please impl RegisterCloudJob", jobName)

		consumerName := jobOrConsumerName
		for _, cfg := range consumers {
			for _, sub := range cfg.Subjects {
				name := c.subjectName(lo.FromPtr(sub.Name))
				assert.If(!allEventKeysSet.Contains(name), "subject:%s not found, please check protobuf define and service", name)
			}

			consumerName = c.consumerName(lo.Ternary(cfg.Consumer != nil, lo.FromPtr(cfg.Consumer), consumerName))
			streamName := c.streamName(cfg.Stream)

			typex.DoBlock(func() {
				if c.consumers[streamName] == nil {
					c.consumers[streamName] = make(map[string]*Consumer)
				}
				assert.If(c.consumers[streamName][consumerName] != nil, "consumer %s already exists", consumerName)

				metadata := map[string]string{"version": fmt.Sprintf("%s/%s", version.Project(), version.Version())}
				consumerCfg := jetstream.ConsumerConfig{
					Name:     consumerName,
					Durable:  consumerName,
					Metadata: metadata,
					AckWait:  time.Minute * 5,
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
					assert.If(c.jobManagers[jobName].managers[subName] == nil, "job manager not found, job=%s subject=%s", jobName, subName)

					job := &jobEventHandler{
						name:         jobName,
						manager:      c.jobManagers[jobName].managers[subName],
						cfg:          subCfg,
						interceptors: c.jobManagers[jobName].interceptors,
					}

					logger.Info().Func(func(e *zerolog.Event) {
						e.Str("job_name", job.name)
						e.Str("job_handler", stack.CallerWithFunc(job.manager.handler).String())
						e.Any("job_config", subCfg)
						e.Any("stream_name", streamName)
						e.Any("consumer_name", consumerName)
						e.Any("job_subject", subName)
						e.Msg("register cloud job manager executor")
					})
					c.jobs[streamName][consumerName][subName] = job
				}
			})
		}
	}
	return
}

func (c *Client) doConsume() (r result.Error) {
	defer result.Recovery(&r)
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
				return r.WithErrorf("concurrent must be in the range of %d-%d", DefaultMinConcurrent, DefaultMaxConcurrent)
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
			c.p.Lc.BeforeStop(lifecycle.WrapNoCtxErr(con.Stop))
		}
	}
	return
}

func (c *Client) Start() error {
	assert.Exit(c.initStream().Err())
	assert.Exit(c.initConsumer().Err())
	assert.Exit(c.doConsume().Err())
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
