package cloudevent

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	yaml "gopkg.in/yaml.v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/typex"
)

const (
	DefaultPrefix             = "acj"
	DefaultTimeout            = 15 * time.Second
	DefaultMaxRetry           = 3
	DefaultRetryBackoff       = time.Second
	DefaultSenderKey          = "sender"
	DefaultCloudEventDelayKey = "__cloudevent_delay_run_at"
	DefaultJobName            = "default"
	DefaultConcurrent         = 100
	DefaultMaxConcurrent      = 1000
	DefaultMinConcurrent      = 1
)

var senderValue = fmt.Sprintf("%s/%s", version.Project(), version.Version())

type Config struct {
	// Streams: nats stream config
	Streams map[string]*StreamConfig `yaml:"streams"`

	// Consumers: nats consumer config
	Consumers map[string]typex.YamlListType[*ConsumerConfig] `yaml:"consumers"`
}

type StreamConfig struct {
	// Storage jetstream.StorageType
	Storage string `yaml:"storage"`

	// Subjects stream subscribe subject, e.g. nvr.speaker.* without prefix
	Subjects typex.YamlListType[string] `yaml:"subjects"`
}

type ConsumerConfig struct {
	// Consumer name without prefix
	Consumer *string `yaml:"consumer"`

	// Concurrent default: 100
	Concurrent *int `yaml:"concurrent"`

	// Stream name without prefix
	Stream string `yaml:"stream"`

	// Subjects config
	Subjects typex.YamlListType[*strOrJobConfig] `yaml:"subjects"`

	// Job event config
	Job *JobEventConfig `yaml:"job"`
}

type JobEventConfig struct {
	// Name subject name
	Name *string `yaml:"name"`

	// Timeout job executor timeout, default: DefaultTimeout
	Timeout *time.Duration `yaml:"timeout"`

	// MaxRetry max retries, default: DefaultMaxRetry
	MaxRetry *int `yaml:"max_retries"`

	// RetryBackoff retry backoff, default: DefaultRetryBackoff
	RetryBackoff *time.Duration `yaml:"retry_backoff"`
}

type jobEventHandler struct {
	// job name
	name string

	// job handler
	handler func(ctx context.Context, args proto.Message) error

	// job config
	cfg *JobEventConfig
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
