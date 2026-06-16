package cloudevent

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pubgo/funk/v2/log"
	cloudeventpb "github.com/pubgo/funk/v2/proto/cloudevent"
	"github.com/pubgo/funk/v2/result"
)

var logger = log.GetLogger("cloudevent")

type Register interface {
	RegisterCloudEvent(jobCli *Client)
}

type RegisterJobOptions struct {
	Opts         *cloudeventpb.RegisterJobOptions
	Interceptors []SubInterceptor
}

type SubInterceptFunc func(ctx context.Context, args proto.Message, handler func(ctx context.Context, args proto.Message) error) error
type SubInterceptor func(next SubInterceptFunc) SubInterceptFunc

type RegisterOpt func(opts *RegisterJobOptions)

type Handler[T proto.Message] func(ctx context.Context, args T) error
type RpcHandler[T proto.Message] func(ctx context.Context, args T) (*emptypb.Empty, error)

type PubOpt func(opts *PubOptions)

type PubInterceptFunc func(ctx context.Context, topic string, args proto.Message, opts *PubOptions, handler func(ctx context.Context, topic string, args proto.Message, opts *PubOptions) result.Result[*PubAckInfo]) result.Result[*PubAckInfo]

type PubInterceptor func(next PubInterceptFunc) PubInterceptFunc

type PubOptions = cloudeventpb.PushEventOptions

type Consumer struct {
	jetstream.Consumer
	Config *ConsumerConfig
}

type PubAckInfo struct {
	AckInfo *jetstream.PubAck
	Header  nats.Header
	MsgId   string
}

type jobManager struct {
	managers     map[string]*handlerManager
	interceptors []SubInterceptor
}

type handlerManager struct {
	handler Handler[proto.Message]
}
