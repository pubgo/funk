package cloudevent

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"

	"github.com/pubgo/funk/v2/assert"
	cloudeventpb "github.com/pubgo/funk/v2/proto/cloudevent"
	"github.com/pubgo/funk/v2/stack"
	"github.com/pubgo/funk/v2/vars"
)

func WrapHandler[Req, Rsp proto.Message](handler func(ctx context.Context, req Req) (Rsp, error)) func(ctx context.Context, req Req) error {
	return func(ctx context.Context, req Req) error {
		_, err := handler(ctx, req)
		return err
	}
}

func init() {
	vars.Register("cloudevent.default_config", func() any {
		return map[string]any{
			"default_prefix":        DefaultPrefix,
			"default_timeout":       DefaultTimeout,
			"default_max_retry":     DefaultMaxRetry,
			"default_retry_backoff": DefaultRetryBackoff,
			"default_job_name":      DefaultJobName,
			"DelayHeaderKey":        DelayHeaderKey,
			"SenderHeaderKey":       SenderHeaderKey,
		}
	})
}

func RegisterJobHandler[T proto.Message](jobCli *Client, jobName, topic string, handler Handler[T], opts ...RegisterOpt) {
	if jobName == "" {
		jobName = DefaultJobName
	}

	eventHandler := func(ctx context.Context, args proto.Message) error { return handler(ctx, args.(T)) }
	jobCli.registerJobHandler(jobName, topic, eventHandler, opts...)
}

func (c *Client) registerJobHandler(jobName, topic string, handler Handler[proto.Message], opts ...RegisterOpt) {
	assert.If(handler == nil, "job manager is nil")
	assert.If(c.subjects[topic] == nil, "topic:%s not found", topic)

	jobOpt := &RegisterJobOptions{Opts: new(cloudeventpb.RegisterJobOptions)}
	for _, o := range opts {
		o(jobOpt)
	}

	if name := lo.FromPtr(jobOpt.Opts.JobName); name != "" {
		jobName = name
	}

	if c.jobManagers[jobName] == nil {
		c.jobManagers[jobName] = &jobManager{managers: make(map[string]*handlerManager), interceptors: jobOpt.Interceptors}
	} else if len(jobOpt.Interceptors) > 0 {
		c.jobManagers[jobName].interceptors = append(c.jobManagers[jobName].interceptors, jobOpt.Interceptors...)
	}

	topic = c.subjectName(topic)
	assert.If(c.jobManagers[jobName].managers[topic] != nil, "job manager already registered, job_name=%s, topic=%s", jobName, topic)

	c.jobManagers[jobName].managers[topic] = &handlerManager{handler: handler}

	logger.Info().Func(func(e *zerolog.Event) {
		e.Str("job", jobName)
		e.Str("topic", topic)
		e.Str("handler", stack.CallerWithFunc(handler).String())
		e.Msg("register cloudevent handler")
	})
}
