package cloudjobs

import (
	"fmt"
	"reflect"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/pkg/gen/cloudjobpb"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/typex"
	"github.com/pubgo/funk/vars"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func (c *Client) RegisterVars() {
	vars.Register("cloudjobs_"+xid.New().String(),
		func() any {
			subjectAndPbType := typex.DoBlock1(func() map[string]string {
				subjectAndPbType := make(map[string]string)
				for k, v := range subjects {
					subjectAndPbType[k] = string(v.ProtoReflect().Descriptor().FullName())
				}
				return subjectAndPbType
			})

			return map[string]any{
				"all_register_subjects":  subjectAndPbType,
				"default_prefix":         DefaultPrefix,
				"default_timeout":        DefaultTimeout,
				"default_max_retry":      DefaultMaxRetry,
				"default_retry_backoff":  DefaultRetryBackoff,
				"default_job_name":       defaultJobName,
				"cloud_job_delay_key":    cloudJobDelayKey,
				"default_concurrent":     defaultConcurrent,
				"default_max_concurrent": defaultMaxConcurrent,
				"default_min_concurrent": defaultMinConcurrent,
				"config":                 c.p.Cfg,
				"prefix":                 c.prefix,
				"streams":                lo.Keys(c.streams),
				"handlers": lo.MapValues(c.handlers,
					func(value map[string]JobHandler[proto.Message], key string) map[string]string {
						return lo.MapValues(value, func(value JobHandler[proto.Message], key string) string {
							return stack.CallerWithFunc(value).String()
						})
					}),
				// jobs map[string]map[string]map[string]*jobHandler
				"jobs": lo.MapValues(c.jobs,
					func(value map[string]map[string]*jobHandler, key string) map[string]map[string]any {
						return lo.MapValues(value, func(value map[string]*jobHandler, key string) map[string]any {
							return lo.MapValues(value, func(value *jobHandler, key string) any {
								return map[string]any{
									"name":    value.name,
									"config":  value.cfg,
									"handler": stack.CallerWithFunc(value.handler).String(),
								}
							})
						})
					}),
				"consumers": lo.MapValues(c.consumers,
					func(value map[string]*Consumer, key string) []string { return lo.Keys(value) }),
			}
		})
}

var subjects = make(map[string]proto.Message)

func RegisterSubject(subject string, subType proto.Message) any {
	assert.If(subject == "", "subject is empty")
	assert.If(subType == nil, "subType is nil")
	assert.If(subjects[subject] != nil, "subject %s already registered", subject)
	logger.Info().Func(func(e *zerolog.Event) {
		e.Str("subject", subject)
		e.Str("type", string(subType.ProtoReflect().Descriptor().FullName()))
		e.Msg("register subject")
	})

	subjects[subject] = subType
	return nil
}

func RegisterJobHandler[T proto.Message](jobCli *Client, jobName string, topic string, handler JobHandler[T], opts ...*cloudjobpb.RegisterJobOptions) {
	assert.Fn(reflect.TypeOf(subjects[topic]) != reflect.TypeOf(lo.Empty[T]()), func() error {
		return fmt.Errorf("type not match, topic-type=%s handler-input-type=%s", reflect.TypeOf(subjects[topic]).String(), reflect.TypeOf(lo.Empty[T]()).String())
	})

	if jobName == "" {
		jobName = defaultJobName
	}

	jobCli.registerJobHandler(jobName, topic, func(ctx *Context, args proto.Message) error { return handler(ctx, args.(T)) }, opts...)
}

func (c *Client) registerJobHandler(jobName string, topic string, handler JobHandler[proto.Message], opts ...*cloudjobpb.RegisterJobOptions) {
	assert.If(handler == nil, "job handler is nil")
	assert.If(subjects[topic] == nil, "topic:%s not found", topic)

	var evtOpt = new(cloudjobpb.RegisterJobOptions)
	for _, o := range opts {
		proto.Merge(evtOpt, o)
	}

	if lo.FromPtr(evtOpt.JobName) != "" {
		jobName = lo.FromPtr(evtOpt.JobName)
	}

	if c.handlers[jobName] == nil {
		c.handlers[jobName] = map[string]JobHandler[proto.Message]{}
	}

	topic = c.subjectName(topic)
	c.handlers[jobName][topic] = handler

	logger.Info().Func(func(e *zerolog.Event) {
		e.Str("job_name", jobName)
		e.Str("topic", topic)
		e.Str("job_handler", stack.CallerWithFunc(handler).String())
		e.Msg("register cloud job handler")
	})
}
