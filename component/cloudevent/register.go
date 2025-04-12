package cloudevent

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pubgo/funk/assert"
	cloudeventpb "github.com/pubgo/funk/proto/cloudevent"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/vars"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func WrapHandler[Req proto.Message, Rsp proto.Message](handler func(ctx context.Context, req Req) (Rsp, error)) func(ctx context.Context, req Req) error {
	return func(ctx context.Context, req Req) error {
		_, err := handler(ctx, req)
		return err
	}
}

func init() {
	vars.Register("cloudevent.default_config", func() any {
		return map[string]any{
			"default_prefix":            DefaultPrefix,
			"default_timeout":           DefaultTimeout,
			"default_max_retry":         DefaultMaxRetry,
			"default_retry_backoff":     DefaultRetryBackoff,
			"default_job_name":          DefaultJobName,
			"DefaultCloudEventDelayKey": DefaultCloudEventDelayKey,
		}
	})
}

func RegisterJobHandler[T proto.Message](jobCli *Client, jobName string, topic string, handler EventHandler[T], opts ...*cloudeventpb.RegisterJobOptions) {
	assert.Fn(reflect.TypeOf(jobCli.subjects[topic]) != reflect.TypeOf(lo.Empty[T]()), func() error {
		return fmt.Errorf("type not match, topic-type=%s handler-input-type=%s", reflect.TypeOf(jobCli.subjects[topic]).String(), reflect.TypeOf(lo.Empty[T]()).String())
	})

	if jobName == "" {
		jobName = DefaultJobName
	}

	jobCli.registerJobHandler(jobName, topic, func(ctx context.Context, args proto.Message) error { return handler(ctx, args.(T)) }, opts...)
}

func (c *Client) registerJobHandler(jobName string, topic string, handler EventHandler[proto.Message], opts ...*cloudeventpb.RegisterJobOptions) {
	assert.If(handler == nil, "job handler is nil")
	assert.If(c.subjects[topic] == nil, "topic:%s not found", topic)

	var evtOpt = new(cloudeventpb.RegisterJobOptions)
	for _, o := range opts {
		proto.Merge(evtOpt, o)
	}

	if lo.FromPtr(evtOpt.JobName) != "" {
		jobName = lo.FromPtr(evtOpt.JobName)
	}

	if c.handlers[jobName] == nil {
		c.handlers[jobName] = map[string]EventHandler[proto.Message]{}
	}

	topic = c.subjectName(topic)
	assert.If(c.handlers[jobName][topic] != nil, "job handler already registered, job_name=%s, topic=%s", jobName, topic)

	c.handlers[jobName][topic] = handler

	logger.Info().Func(func(e *zerolog.Event) {
		e.Str("job_name", jobName)
		e.Str("topic", topic)
		e.Str("job_handler", stack.CallerWithFunc(handler).String())
		e.Msg("register cloud job handler")
	})
}
