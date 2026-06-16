package cloudevent

import (
	"context"
	"net/http"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/component/lifecycle"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/stack"
	"github.com/pubgo/funk/v2/try"
)

func (c *Client) doConsumeHandler(streamName, consumerName string, jobSubjects map[string]*jobEventHandler, concurrent int) func(msg jetstream.Msg) {
	handler := func(msg jetstream.Msg) {
		now := time.Now()
		addMsgInfo := func(e *zerolog.Event) {
			e.Str("stream", streamName)
			e.Str("consumer", consumerName)
			e.Any("header", msg.Headers())
			e.Any("msg_id", msg.Headers().Get(jetstream.MsgIDHeader))
			e.Str("subject", msg.Subject())
			e.Str("msg_received_time", now.String())
			e.Str("job_cost", time.Since(now).String())
		}

		logger.Debug().Func(addMsgInfo).Msg("received cloud job manager")

		handlerDelayJob := func() (r result.Result[bool]) {
			defer result.Recovery(&r)
			dur := decodeDelayTime(msg.Headers().Get(DelayHeaderKey)).
				MapErr(func(err error) error {
					return errors.Wrap(err, "failed to parse job delay time")
				}).
				UnwrapOrThrow(&r)
			if r.IsErr() {
				return
			}

			if dur <= 0 {
				return r.WithValue(false)
			}

			if result.Throw(&r, msg.NakWithDelay(dur)) {
				return
			}

			return r.WithValue(true)
		}

		delayRet := handlerDelayJob().
			IfErr(func(err error) {
				logger.Err(err).Func(addMsgInfo).Msg("failed to handle cloud delay job and no ack")
			}).
			IfOK(func(b bool) {
				if b {
					logger.Info().Func(addMsgInfo).Msg("redeliver the message after the given delay")
				}
			})

		if delayRet.IsErr() || delayRet.Unwrap() {
			return
		}

		job := jobSubjects[msg.Subject()]
		if job == nil {
			logger.Error().Func(addMsgInfo).Msg("failed to find subject job manager")
			return
		}

		meta := result.Wrap(msg.Metadata()).
			IfErr(func(err error) {
				logger.Err(err).Func(addMsgInfo).Msg("failed to parse nats stream msg metadata")
			})
		if meta.IsErr() {
			return
		}

		cfg := job.cfg
		checkErrAndLog := func(err error, msg string) {
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

		err := try.Try(func() error { return c.doHandler(meta.Unwrap(), msg, job, cfg).Err() })
		if err == nil {
			checkErrAndLog(msg.Ack(), "failed to do msg ack with manager ok")
			return
		}

		if isRejectErr(err) {
			checkErrAndLog(msg.TermWithReason("reject by caller"), "failed to do msg ack with reject err")
			return
		}

		if isForceRetry(err) {
			checkErrAndLog(msg.Nak(), "force_retry: failed to reply nak msg")
			return
		}

		backoff := lo.FromPtr(cfg.RetryBackoff)
		maxRetries := lo.FromPtr(cfg.MaxRetry)

		if err1 := isRedeliveryErr(err); err1 != nil {
			backoff = err1.delay
		}

		if meta.Unwrap().NumDelivered < uint64(maxRetries) {
			logger.Warn().
				Err(err).
				Func(addMsgInfo).
				Any("metadata", meta).
				Msg("retry nats stream cloud job manager")
			checkErrAndLog(msg.NakWithDelay(backoff), "failed to retry msg with delay nak")
			return
		}

		checkErrAndLog(err, "failed to do manager cloud job")
		checkErrAndLog(msg.Ack(), "failed to do msg ack with manager error")
	}

	pool := assert.Must1(ants.NewPool(
		concurrent,
		ants.WithLogger(log.NewStd(logger)),
		ants.WithNonblocking(false),
	))
	c.p.Lc.BeforeStop(lifecycle.WrapNoCtxErr(func() {
		pool.Release()
	}))
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

func (c *Client) doHandler(meta *jetstream.MsgMetadata, msg jetstream.Msg, job *jobEventHandler, cfg *JobEventConfig) (gErr result.Error) {
	defer result.Recovery(&gErr)
	timeout := lo.FromPtr(cfg.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx = log.UpdateFieldsCtx(ctx, log.Fields{
		"sub_subject":   msg.Subject(),
		"sub_stream":    meta.Stream,
		"sub_consumer":  meta.Consumer,
		"sub_msg_id":    msg.Headers().Get(jetstream.MsgIDHeader),
		SenderHeaderKey: msg.Headers().Get(SenderHeaderKey),
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

	now := time.Now()
	var args any
	defer gErr.InspectErr(func(err error) {
		logger.Err(err).Func(func(e *zerolog.Event) {
			e.Any("context", msgCtx)
			e.Any("args", args)
			e.Str("timeout", timeout.String())
			e.Str("start_time", now.String())
			e.Str("job_cost", time.Since(now).String())
			e.Msg("failed to do cloud job manager")
		})
	})

	var pb anypb.Any
	if result.ErrOf(proto.Unmarshal(msg.Data(), &pb)).
		MapErr(func(err error) error {
			return errors.WrapTags(err, errors.Tags{
				"msg":  "failed to unmarshal stream msg data to any proto",
				"args": string(msg.Data()),
			})
		}).
		Throw(&gErr) {
		return
	}
	args = &pb

	dst := result.Wrap(anypb.UnmarshalNew(args.(*anypb.Any), proto.UnmarshalOptions{})).
		MapErr(func(err error) error {
			return errors.WrapTags(err, errors.Tags{
				"msg":  "failed to unmarshal any proto to proto msg",
				"args": args,
			})
		}).
		UnwrapOrThrow(&gErr)
	if gErr.IsErr() {
		return
	}

	ctx = createCtxWithSubjectContext(ctx, msgCtx)

	interceptor := func(ctx context.Context, args proto.Message, handler func(ctx context.Context, args proto.Message) error) error {
		return handler(ctx, args)
	}

	for i := len(job.interceptors) - 1; i >= 0; i-- {
		interceptor = job.interceptors[i](interceptor)
	}

	return result.ErrOf(interceptor(ctx, dst, job.manager.handler)).
		MapErr(func(err error) error {
			return errors.WrapTags(err, errors.Tags{
				"msg":    "failed to do cloud job manager",
				"args":   args,
				"any_pb": dst,
			})
		})
}
