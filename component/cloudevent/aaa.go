package cloudevent

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pubgo/catdogs/pkg/gen/proto/cloudeventpb"
	"github.com/pubgo/funk/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

var logger = log.GetLogger("cloudevent")

type EventRegister interface {
	RegisterCloudEvent(jobCli *Client)
}

type EventHandler[T proto.Message] func(ctx context.Context, args T) error
type RpcEventHandler[T proto.Message] func(ctx context.Context, args T) (*emptypb.Empty, error)

type Options = cloudeventpb.PushEventOptions

type PushEventOpt func(opts *Options)

type Consumer struct {
	jetstream.Consumer
	Config *ConsumerConfig
}

type PubAckInfo struct {
	AckInfo *jetstream.PubAck
	Header  nats.Header
	MsgId   string
}
