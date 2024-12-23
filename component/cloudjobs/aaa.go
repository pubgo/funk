package cloudjobs

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/pkg/gen/cloudjobpb"
	"google.golang.org/protobuf/proto"
)

var logger = log.GetLogger("cloudjobs")

type Register interface {
	RegisterCloudJobs(jobCli *Client)
}

type JobHandler[T proto.Message] func(ctx *Context, args T) error

type Options = cloudjobpb.PushEventOptions

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
