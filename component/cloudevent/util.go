package cloudevent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	cloudeventpb "github.com/pubgo/funk/proto/cloudevent"
	"github.com/pubgo/funk/protoutils"
	"github.com/pubgo/funk/result"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func getStorageType(name string) jetstream.StorageType {
	switch name {
	case "memory":
		return jetstream.MemoryStorage
	case "file":
		return jetstream.FileStorage
	default:
		panic("unknown storage type")
	}
}

func mergeJobConfig(dst *JobEventConfig, src *JobEventConfig) *JobEventConfig {
	if src == nil {
		src = handleDefaultJobConfig(nil)
	}

	if dst == nil {
		dst = handleDefaultJobConfig(nil)
	}

	if dst.MaxRetry == nil {
		dst.MaxRetry = src.MaxRetry
	}

	if dst.Timeout == nil {
		dst.Timeout = src.Timeout
	}

	if dst.RetryBackoff == nil {
		dst.RetryBackoff = src.RetryBackoff
	}

	return dst
}

func handleDefaultJobConfig(cfg *JobEventConfig) *JobEventConfig {
	if cfg == nil {
		cfg = new(JobEventConfig)
	}

	if cfg.Timeout == nil {
		cfg.Timeout = lo.ToPtr(DefaultTimeout)
	}

	if cfg.MaxRetry == nil {
		cfg.MaxRetry = lo.ToPtr(DefaultMaxRetry)
	}

	if cfg.RetryBackoff == nil {
		cfg.RetryBackoff = lo.ToPtr(DefaultRetryBackoff)
	}

	return cfg
}

func handleSubjectName(name string, prefix string) string {
	prefix = fmt.Sprintf("%s.", prefix)
	if strings.HasPrefix(name, prefix) {
		return name
	}

	return fmt.Sprintf("%s%s", prefix, name)
}

func encodeDelayTime(duration time.Duration) string {
	return strconv.Itoa(int(time.Now().Add(duration).UnixMilli()))
}

func decodeDelayTime(delayTime string) (r result.Result[time.Duration]) {
	tt := result.Wrap(strconv.Atoi(delayTime)).
		MapErr(func(err error) error {
			return errors.Wrapf(err, "failed to parse cloud event job delay time, time=%s", delayTime)
		})

	return result.MapTo(tt, func(t int) time.Duration {
		return time.Until(time.UnixMilli(int64(tt.GetValue())))
	})
}

type subjectOpt struct {
	*cloudeventpb.CloudEventServiceOptions
	*cloudeventpb.CloudEventMethodOptions
}

func registerSubject(subjects map[string]*cloudeventpb.CloudEventMethodOptions, subject string, operation string, data *cloudeventpb.CloudEventMethodOptions) any {
	assert.If(subject == "", "subject is empty")
	assert.If(operation == "", "operation is empty")
	assert.If(data == nil, "data is nil")
	assert.If(subjects[subject] != nil, "subject %s already registered", subject)
	assert.If(subjects[operation] != nil, "operation %s already registered", operation)
	logger.Info().Func(func(e *zerolog.Event) {
		e.Str("subject", subject)
		e.Str("operation", operation)
		e.Any("details", data)
		e.Msg("register subject")
	})

	data.Operation = lo.ToPtr(operation)
	subjects[subject] = data
	subjects[operation] = data
	return nil
}

func getAllSubject() map[string]*cloudeventpb.CloudEventMethodOptions {
	var subjects = make(map[string]*cloudeventpb.CloudEventMethodOptions)
	for _, opt := range getAllSubjectOptions() {
		registerSubject(subjects, opt.CloudEventServiceOptions.Name, *opt.Operation, opt.CloudEventMethodOptions)
	}
	return subjects
}

func getAllSubjectOptions() []subjectOpt {
	var opts []subjectOpt
	protoutils.EachService(func(desc protoreflect.FileDescriptor, srv protoreflect.ServiceDescriptor) {
		if !protoutils.HasExtension(srv.Options(), cloudeventpb.E_Job) {
			return
		}

		jobOpt := protoutils.GetExtension[cloudeventpb.CloudEventServiceOptions](srv.Options(), cloudeventpb.E_Job)
		if jobOpt == nil {
			return
		}

		protoutils.EachServiceMethod(srv, func(mth protoreflect.MethodDescriptor) {
			if !protoutils.HasExtension(mth.Options(), cloudeventpb.E_Subject) {
				return
			}

			subOpt := protoutils.GetExtension[cloudeventpb.CloudEventMethodOptions](mth.Options(), cloudeventpb.E_Subject)
			if subOpt == nil {
				return
			}

			subOpt.Operation = lo.ToPtr(fmt.Sprintf("/%s/%s", srv.FullName(), mth.Name()))
			opts = append(opts, subjectOpt{CloudEventServiceOptions: jobOpt, CloudEventMethodOptions: subOpt})
		})
	})
	return opts
}
