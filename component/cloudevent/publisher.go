package cloudevent

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/pubgo/funk/v2/ctxutil"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/typex"
)

func Publish(jobCli *Client, ctx context.Context, topic string, args proto.Message, interceptors []PubInterceptor, opts ...PubOpt) result.Result[*PubAckInfo] {
	return jobCli.Publish(ctx, topic, args, interceptors, opts...)
}

func (c *Client) Publish(ctx context.Context, topic string, args proto.Message, interceptors []PubInterceptor, opts ...PubOpt) result.Result[*PubAckInfo] {
	return c.publish(ctx, topic, args, interceptors, opts...)
}

func (c *Client) doPublish(ctx context.Context, topic string, args proto.Message, opts *PubOptions) (r result.Result[*PubAckInfo]) {
	defer result.Recovery(&r)

	if opts == nil {
		opts = new(PubOptions)
	}

	msgId := xid.New().String()
	if opts.MsgId != nil {
		msgId = opts.GetMsgId()
	}

	pb := result.Wrap(anypb.New(args)).
		MapErr(func(err error) error {
			return errors.Wrap(err, "failed to marshal args to any proto")
		}).
		UnwrapOrThrow(&r)
	if r.IsErr() {
		return
	}

	data := result.Wrap(proto.Marshal(pb)).
		MapErr(func(err error) error {
			return errors.Wrap(err, "failed to marshal any proto to bytes")
		}).
		UnwrapOrThrow(&r)
	if r.IsErr() {
		return
	}

	topic = c.subjectName(topic)
	header := typex.DoBlock1(func() nats.Header {
		header := nats.Header{
			SenderHeaderKey: []string{lo.FromPtr(opts.Sender)},
			DelayHeaderKey:  []string{encodeDelayTime(opts.Delay)},
		}
		for k, v := range opts.Metadata {
			header.Add(k, v)
		}
		return header
	})

	msg := &nats.Msg{Subject: topic, Data: data, Header: header}
	jetOpts := append([]jetstream.PublishOpt{}, jetstream.WithMsgID(msgId))
	pubActInfo := result.Wrap(c.js.PublishMsg(ctx, msg, jetOpts...)).
		MapErr(func(err error) error {
			return errors.Wrapf(err, "failed to publish msg to jetstream, topic=%s msg_id=%s", topic, msgId)
		}).
		UnwrapOrThrow(&r)
	if r.IsErr() {
		return
	}

	return r.WithValue(&PubAckInfo{
		AckInfo: pubActInfo,
		Header:  header,
		MsgId:   msgId,
	})
}

func (c *Client) publish(ctx context.Context, topic string, args proto.Message, interceptors []PubInterceptor, opts ...PubOpt) (r result.Result[*PubAckInfo]) {
	timeout := ctxutil.GetTimeout(ctx)
	now := time.Now()
	logFn := func(e *zerolog.Event) {
		e.Str("topic", topic)
		e.Str("start_at", now.String())
		e.Any("args", args)
		e.Str("cost", time.Since(now).String())
		e.Any("ack_info", r.UnwrapOrEmpty())
		if timeout != nil {
			e.Str("timeout", timeout.String())
		}
	}

	defer func() {
		if r.IsOK() {
			logger.Info(ctx).Func(logFn).Msg("succeed to publish cloudevent msg to jetstream")
		} else {
			logger.Err(r.Err(), ctx).Func(logFn).Msg("failed to publish cloudevent msg to jetstream")
		}
	}()

	interceptor := func(ctx context.Context, topic string, args proto.Message, opts *PubOptions, handler func(ctx context.Context, topic string, args proto.Message, opts *PubOptions) result.Result[*PubAckInfo]) result.Result[*PubAckInfo] {
		return handler(ctx, topic, args, opts)
	}
	for i := len(interceptors) - 1; i >= 0; i-- {
		interceptor = interceptors[i](interceptor)
	}

	pushEventOpt := getOptions(ctx, opts...)
	return interceptor(ctx, topic, args, pushEventOpt, c.doPublish)
}
