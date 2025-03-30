package cloudevent

import (
	"context"
	"net/http"
	"time"

	"github.com/pubgo/catdogs/pkg/gen/proto/cloudeventpb"
	"github.com/rs/xid"
	"google.golang.org/protobuf/proto"
)

type Context struct {
	// Header jetstream.Headers().
	Header http.Header

	// NumDelivered jetstream.MsgMetadata{}.NumDelivered
	NumDelivered uint64

	// NumPending jetstream.MsgMetadata{}.NumPending
	NumPending uint64

	// Timestamp jetstream.MsgMetadata{}.Timestamp
	Timestamp time.Time

	// Stream jetstream.MsgMetadata{}.Stream
	Stream string

	// Consumer jetstream.MsgMetadata{}.Consumer
	Consumer string

	// Subject|Topic name jetstream.Msg().Subject()
	Subject string

	// Config job config from config file or default
	Config *JobEventConfig
}

var cloudeventCtxKey = xid.New().String()

func createCtxWithContext(parent context.Context, ctx *Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, cloudeventCtxKey, ctx)
}

func GetEventContext(ctx context.Context) *Context {
	if ctx == nil {
		return nil
	}

	evtCtx, ok := ctx.Value(cloudeventCtxKey).(*Context)
	if !ok {
		return nil
	}

	return evtCtx
}

var pushEventCtxKey = xid.New().String()

func withOptions(ctx context.Context, opts ...*cloudeventpb.PushEventOptions) context.Context {
	if len(opts) == 0 {
		return ctx
	}

	oldOpts, ok := ctx.Value(pushEventCtxKey).(*cloudeventpb.PushEventOptions)
	if !ok {
		oldOpts = new(cloudeventpb.PushEventOptions)
	}

	for i := range opts {
		proto.Merge(oldOpts, opts[i])
	}

	return context.WithValue(ctx, pushEventCtxKey, oldOpts)
}

func WithPushOpt(opts ...func(opt *cloudeventpb.PushEventOptions)) *cloudeventpb.PushEventOptions {
	var opt cloudeventpb.PushEventOptions
	for _, o := range opts {
		o(&opt)
	}
	return &opt
}

func getOptions(ctx context.Context, opts ...*cloudeventpb.PushEventOptions) *cloudeventpb.PushEventOptions {
	var evtOpt = new(cloudeventpb.PushEventOptions)
	opt, ok := ctx.Value(pushEventCtxKey).(*cloudeventpb.PushEventOptions)
	if ok {
		evtOpt = opt
	}

	for _, o := range opts {
		proto.Merge(evtOpt, o)
	}

	if evtOpt.GetMsgId() == "" {
		evtOpt.MsgId = nil
	}

	return evtOpt
}
