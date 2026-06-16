package cloudevent

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	yaml "gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/typex"
)

const (
	DefaultPrefix        = "acj"
	DefaultTimeout       = 15 * time.Second
	DefaultMaxRetry      = 3
	DefaultRetryBackoff  = time.Second
	SenderHeaderKey      = "__cloudevent_sender"
	DelayHeaderKey       = "__cloudevent_delay_run_at"
	DefaultJobName       = "default"
	DefaultConcurrent    = 100
	DefaultMaxConcurrent = 1000
	DefaultMinConcurrent = 1
)

var senderValue = fmt.Sprintf("%s/%s", version.Project(), version.Version())

type Config struct {
	Streams   map[string]*StreamConfig                       `yaml:"streams"`
	Consumers map[string]typex.YamlListType[*ConsumerConfig] `yaml:"consumers"`
}

type StreamConfig struct {
	Storage  string                     `yaml:"storage"`
	Subjects typex.YamlListType[string] `yaml:"subjects"`
}

type ConsumerConfig struct {
	Consumer   *string                             `yaml:"consumer"`
	Concurrent *int                                `yaml:"concurrent"`
	Stream     string                              `yaml:"stream"`
	Subjects   typex.YamlListType[*strOrJobConfig] `yaml:"subjects"`
	Job        *JobEventConfig                     `yaml:"job"`
}

type JobEventConfig struct {
	Name         *string        `yaml:"name"`
	Timeout      *time.Duration `yaml:"timeout"`
	MaxRetry     *int           `yaml:"max_retries"`
	RetryBackoff *time.Duration `yaml:"retry_backoff"`
}

type jobEventHandler struct {
	name         string
	manager      *handlerManager
	cfg          *JobEventConfig
	interceptors []SubInterceptor
}

func consumerAckWait(cfg *ConsumerConfig) time.Duration {
	const minAckWait = 5 * time.Minute

	base := handleDefaultJobConfig(cfg.Job)
	maxTimeout := lo.FromPtr(base.Timeout)
	for _, sub := range cfg.Subjects {
		subCfg := mergeJobConfig(lo.ToPtr(JobEventConfig(lo.FromPtr(sub))), base)
		if t := lo.FromPtr(subCfg.Timeout); t > maxTimeout {
			maxTimeout = t
		}
	}

	ackWait := maxTimeout + time.Minute
	if ackWait < minAckWait {
		return minAckWait
	}
	return ackWait
}

type strOrJobConfig JobEventConfig

func (p *strOrJobConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.IsZero() {
		return nil
	}

	switch value.Kind {
	case yaml.ScalarNode:
		var data string
		if err := value.Decode(&data); err != nil {
			return errors.WrapCaller(err)
		}

		*p = strOrJobConfig(JobEventConfig{Name: &data})
		return nil
	case yaml.MappingNode:
		var data JobEventConfig
		if err := value.Decode(&data); err != nil {
			return errors.WrapCaller(err)
		}

		*p = strOrJobConfig(data)
		return nil
	default:
		var val any
		assert.Exit(value.Decode(&val))
		return errors.Errorf("yaml kind type error,kind=%v data=%v", value.Kind, val)
	}
}
