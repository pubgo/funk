package pyroscope

import (
	"fmt"

	"github.com/grafana/pyroscope-go"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/component/lifecycle"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/running"
)

type Param struct {
	Cfg    *Config
	Logger log.Logger
	Lc     lifecycle.Lifecycle
}

type Client struct {
	profiler *pyroscope.Profiler
}

func (c *Client) Enabled() bool {
	return c != nil && c.profiler != nil
}

func (c *Client) Profiler() *pyroscope.Profiler {
	if c == nil {
		return nil
	}
	return c.profiler
}

func (c *Client) Flush(wait bool) {
	if c == nil || c.profiler == nil {
		return
	}
	c.profiler.Flush(wait)
}

func New(p Param) *Client {
	cfg := mergeConfig(p.Cfg)
	if !cfg.Enabled || cfg.ServerAddress == "" {
		return &Client{}
	}

	logger := p.Logger
	if logger == nil {
		logger = log.GetLogger(Name)
	}

	pyroCfg := pyroscope.Config{
		ApplicationName:   applicationName(cfg),
		ServerAddress:     cfg.ServerAddress,
		BasicAuthUser:     cfg.BasicAuthUser,
		BasicAuthPassword: cfg.BasicAuthPassword,
		TenantID:          cfg.TenantID,
		UploadRate:        cfg.UploadRate,
		ProfileTypes:      cfg.profileTypes(),
		Tags:              defaultTags(cfg),
		DisableGCRuns:     cfg.DisableGCRuns,
	}
	if !cfg.DisableLog {
		pyroCfg.Logger = newLoggerAdapter(logger)
	}

	profiler := assert.Must1(pyroscope.Start(pyroCfg))
	logger.Info().
		Str("application_name", pyroCfg.ApplicationName).
		Str("server_address", pyroCfg.ServerAddress).
		Msg("pyroscope profiler started")

	if p.Lc != nil {
		p.Lc.BeforeStop(lifecycle.WrapNoCtxErr(func() {
			if err := profiler.Stop(); err != nil {
				logger.Err(err).Msg("failed to stop pyroscope profiler")
			}
		}))
	}

	return &Client{profiler: profiler}
}

func applicationName(cfg *Config) string {
	if cfg.ApplicationName != "" {
		return cfg.ApplicationName
	}
	return fmt.Sprintf("%s/%s", running.Project(), version.Version())
}

func defaultTags(cfg *Config) map[string]string {
	tags := map[string]string{
		"hostname":    running.Hostname,
		"env":         running.Env.String(),
		"version":     running.Version(),
		"instance_id": running.InstanceID,
	}
	for k, v := range cfg.Tags {
		tags[k] = v
	}
	return tags
}
