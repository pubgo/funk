package pyroscope

import (
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/merge"
)

const Name = "pyroscope"

type Config struct {
	Enabled           bool              `yaml:"enabled"`
	ApplicationName   string            `yaml:"application_name"`
	ServerAddress     string            `yaml:"server_address"`
	BasicAuthUser     string            `yaml:"basic_auth_user"`
	BasicAuthPassword string            `yaml:"basic_auth_password"`
	TenantID          string            `yaml:"tenant_id"`
	UploadRate        time.Duration     `yaml:"upload_rate"`
	ProfileTypes      []string          `yaml:"profile_types"`
	Tags              map[string]string `yaml:"tags"`
	DisableGCRuns     bool              `yaml:"disable_gc_runs"`
	DisableLog        bool              `yaml:"disable_log"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:      false,
		ProfileTypes: defaultProfileTypeNames(),
		UploadRate:   time.Minute,
	}
}

func defaultProfileTypeNames() []string {
	return lo.Map(pyroscope.DefaultProfileTypes, func(item pyroscope.ProfileType, _ int) string {
		return string(item)
	})
}

func (c *Config) profileTypes() []pyroscope.ProfileType {
	if len(c.ProfileTypes) == 0 {
		return pyroscope.DefaultProfileTypes
	}

	types := make([]pyroscope.ProfileType, 0, len(c.ProfileTypes))
	for _, name := range c.ProfileTypes {
		types = append(types, pyroscope.ProfileType(name))
	}
	return types
}

func mergeConfig(cfg *Config) *Config {
	if cfg == nil {
		return DefaultConfig()
	}
	return merge.Copy(DefaultConfig(), cfg).Unwrap()
}
