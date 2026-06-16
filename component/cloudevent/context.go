package cloudevent

import (
	"context"
	"net/http"
	"time"

	cloudeventpb "github.com/pubgo/funk/v2/proto/cloudevent"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

type Context struct {
	Header       http.Header
	NumDelivered uint64
	NumPending   uint64
	Timestamp    time.Time
	Stream       string
	Consumer     string
	Subject      string
	Config       *JobEventConfig
}

type ctxKey int

const cloudeventCtxKey ctxKey = 1

func createCtxWithSubjectContext(parent context.Context, ctx *Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, cloudeventCtxKey, ctx)
}

func GetContext(ctx context.Context) *Context {
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

func WithPushOpt(opts ...func(opt *cloudeventpb.PushEventOptions)) *cloudeventpb.PushEventOptions {
	var opt cloudeventpb.PushEventOptions
	for _, o := range opts {
		o(&opt)
	}
	return &opt
}

func getOptions(ctx context.Context, opts ...PubOpt) *PubOptions {
	evtOpt := new(PubOptions)
	if opt, ok := ctx.Value(pushEventCtxKey).(*PubOptions); ok {
		evtOpt = opt
	}

	for _, o := range opts {
		if o == nil {
			continue
		}
		o(evtOpt)
	}

	if evtOpt.GetMsgId() == "" {
		evtOpt.MsgId = nil
	}

	if evtOpt.ContentType == nil {
		evtOpt.ContentType = lo.ToPtr("application/json")
	}

	if evtOpt.Sender == nil {
		evtOpt.Sender = lo.ToPtr(senderValue)
	}

	return evtOpt
}

func ProtoPubOpts(opts ...*cloudeventpb.PushEventOptions) []PubOpt {
	if len(opts) == 0 {
		return nil
	}
	return []PubOpt{func(po *PubOptions) {
		for _, o := range opts {
			if o == nil {
				continue
			}
			proto.Merge(po, o)
		}
	}}
}

func ProtoRegisterOpts(opts ...*cloudeventpb.RegisterJobOptions) RegisterOpt {
	return func(ro *RegisterJobOptions) {
		if ro.Opts == nil {
			ro.Opts = new(cloudeventpb.RegisterJobOptions)
		}
		for _, o := range opts {
			if o == nil {
				continue
			}
			proto.Merge(ro.Opts, o)
		}
	}
}

func WithSubInterceptors(interceptors ...SubInterceptor) RegisterOpt {
	return func(ro *RegisterJobOptions) {
		ro.Interceptors = append(ro.Interceptors, interceptors...)
	}
}
