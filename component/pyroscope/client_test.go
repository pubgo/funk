package pyroscope

import (
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/stretchr/testify/assert"
)

func TestMergeConfigDefaults(t *testing.T) {
	cfg := mergeConfig(&Config{
		Enabled:       true,
		ServerAddress: "http://localhost:4040",
	})
	assert.Equal(t, defaultProfileTypeNames(), cfg.ProfileTypes)
	assert.Equal(t, time.Minute, cfg.UploadRate)
}

func TestMergeConfigNil(t *testing.T) {
	cfg := mergeConfig(nil)
	assert.Equal(t, DefaultConfig(), cfg)
}

func TestProfileTypesFallback(t *testing.T) {
	cfg := mergeConfig(&Config{})
	assert.Equal(t, pyroscope.DefaultProfileTypes, cfg.profileTypes())
}

func TestApplicationNameDefault(t *testing.T) {
	name := applicationName(&Config{})
	assert.NotEmpty(t, name)
}

func TestDefaultTags(t *testing.T) {
	tags := defaultTags(&Config{Tags: map[string]string{"service": "api"}})
	assert.Equal(t, "api", tags["service"])
	assert.NotEmpty(t, tags["hostname"])
	assert.NotEmpty(t, tags["env"])
}

func TestNewDisabled(t *testing.T) {
	client := New(Param{Cfg: &Config{Enabled: false, ServerAddress: "http://localhost:4040"}})
	assert.False(t, client.Enabled())
}

func TestNewMissingServer(t *testing.T) {
	client := New(Param{Cfg: &Config{Enabled: true}})
	assert.False(t, client.Enabled())
}
