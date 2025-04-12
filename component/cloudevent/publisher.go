package cloudevent

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pubgo/funk/ctxutil"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/errors/errcheck"
	cloudeventpb "github.com/pubgo/funk/proto/cloudevent"
	"github.com/pubgo/funk/result"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
	"github.com/pubgo/funk/typex"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func PushEvent[T any](handler func(*Client, context.Context, T, ...*cloudeventpb.PushEventOptions) (*PubAckInfo, error), jobCli *Client, ctx context.Context, t T, opts ...*cloudeventpb.PushEventOptions) chan result.Result[*PubAckInfo] {
	errChan := make(chan result.Result[*PubAckInfo])
	timeout := ctxutil.GetTimeout(ctx)
	now := time.Now()
	fnCaller := stack.Caller(1).String()

	// clone ctx and recalculate timeout
	ctx = lo.T2(ctxutil.Clone(ctx, DefaultTimeout)).A
	go func() {
		var getPubAck = func() (pubAck *PubAckInfo, err error) {
			err = try.Try(func() error {
				pubAck, err = handler(jobCli, ctx, t, opts...)
				return err
			})
			return
		}
		pubAck, err := getPubAck()
		err = errors.IfErr(err, func(err error) error {
			logger.Err(err, ctx).Func(func(e *zerolog.Event) {
				if timeout != nil {
					e.Str("timeout", timeout.String())
				}

				e.Str("fn_caller", fnCaller)
				e.Any("input", t)
				e.Str("stack", stack.CallerWithFunc(handler).String())
				e.Str("cost", time.Since(now).String())
				e.Msg("failed to push event msg to nats stream")
			})
			return err
		})
		errChan <- result.Wrap(pubAck, err)
	}()
	return errChan
}

// PushRpcEvent push event async
func PushRpcEvent[T proto.Message](handler RpcEventHandler[T], ctx context.Context, t T, opts ...*cloudeventpb.PushEventOptions) chan error {
	// clone ctx and recalculate timeout
	ctx = lo.T2(ctxutil.Clone(ctx, DefaultTimeout)).A
	ctx = withOptions(ctx, opts...)

	fnCaller := stack.Caller(1).String()
	errChan := make(chan error)
	timeout := ctxutil.GetTimeout(ctx)
	now := time.Now()

	var pushEventBasic = func(handler RpcEventHandler[T], ctx context.Context) error {
		err := try.Try(func() error { return lo.T2(handler(ctx, t)).B })
		if err == nil {
			return nil
		}

		logger.Err(err, ctx).Func(func(e *zerolog.Event) {
			if timeout != nil {
				e.Str("timeout", timeout.String())
			}

			e.Str("fn_caller", fnCaller)
			e.Any("input", t)
			e.Str("stack", stack.CallerWithFunc(handler).String())
			e.Str("cost", time.Since(now).String())
			e.Msg("failed to push event msg to nats stream")
		})
		return err
	}

	go func() { errChan <- pushEventBasic(handler, ctx) }()
	return errChan
}

func (c *Client) Publish(ctx context.Context, topic string, args proto.Message, opts ...*cloudeventpb.PushEventOptions) (*PubAckInfo, error) {
	return c.publish(ctx, topic, args, opts...)
}

func (c *Client) publish(ctx context.Context, topic string, args proto.Message, opts ...*cloudeventpb.PushEventOptions) (_ *PubAckInfo, gErr error) {
	var timeout = ctxutil.GetTimeout(ctx)
	var now = time.Now()
	var msgId = xid.New().String()
	var pushEventOpt *cloudeventpb.PushEventOptions
	var pubActInfo *jetstream.PubAck

	defer func() {
		var msgFn = func(e *zerolog.Event) {
			e.Str("pub_topic", topic)
			e.Str("pub_start", now.String())
			e.Any("pub_args", args)
			e.Str("pub_cost", time.Since(now).String())
			e.Str("pub_msg_id", msgId)
			e.Any("pub_ack_info", pubActInfo)
			if timeout != nil {
				e.Str("timeout", timeout.String())
			}
		}
		if gErr == nil {
			logger.Info(ctx).Func(msgFn).Msg("succeed to publish cloud event job to stream")
		} else {
			logger.Err(gErr, ctx).Func(msgFn).Msg("failed to publish cloud event job to stream")
		}
	}()

	pushEventOpt = getOptions(ctx, opts...)
	if pushEventOpt.MsgId != nil {
		msgId = pushEventOpt.GetMsgId()
	}

	pb, err := anypb.New(args)
	err = errors.IfErr(err, func(err error) error {
		return errors.Wrap(err, "failed to marshal args to any proto")
	})
	if errcheck.Check(&gErr, err) {
		return
	}

	// TODO get parent event info from ctx
	data, err := proto.Marshal(pb)
	err = errors.IfErr(err, func(err error) error {
		return errors.Wrap(err, "failed to marshal any proto to bytes")
	})
	if errcheck.Check(&gErr, err) {
		return
	}

	// subject|topic name
	topic = c.subjectName(topic)
	header := typex.DoBlock1(func() nats.Header {
		header := nats.Header{
			DefaultSenderKey:          []string{senderValue},
			DefaultCloudEventDelayKey: []string{encodeDelayTime(pushEventOpt.DelayDur.AsDuration())},
		}
		for k, v := range pushEventOpt.Metadata {
			header.Add(k, v)
		}
		return header
	})

	msg := &nats.Msg{Subject: topic, Data: data, Header: header}
	jetOpts := append([]jetstream.PublishOpt{}, jetstream.WithMsgID(msgId))
	pubActInfo, err = c.js.PublishMsg(ctx, msg, jetOpts...)
	err = errors.IfErr(err, func(err error) error {
		return errors.Wrapf(err, "failed to publish msg to stream, topic=%s msg_id=%s", topic, msgId)
	})
	if errcheck.Check(&gErr, err) {
		return
	}

	return &PubAckInfo{
		AckInfo: pubActInfo,
		Header:  header,
		MsgId:   msgId,
	}, nil
}
