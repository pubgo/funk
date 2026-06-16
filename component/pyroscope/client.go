package pyroscope

import (
	"fmt"

	"github.com/grafana/pyroscope-go"
	"github.com/samber/lo"

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
	logger   log.Logger
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
	if c.logger != nil {
		c.logger.Debug().Bool("wait", wait).Msg("flushing pyroscope profiler session")
	}
	c.profiler.Flush(wait)
	if c.logger != nil {
		c.logger.Debug().Bool("wait", wait).Msg("pyroscope profiler session flushed")
	}
}

func New(p Param) *Client {
	cfg := mergeConfig(p.Cfg)
	logger := resolveLogger(p)

	if !cfg.Enabled {
		logger.Info().Msg("pyroscope profiler disabled by config")
		return &Client{logger: logger}
	}
	if cfg.ServerAddress == "" {
		logger.Warn().Msg("pyroscope profiler skipped: server_address is empty")
		return &Client{logger: logger}
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

	logger.Info().Func(func(e *log.Event) {
		e.Str("application_name", pyroCfg.ApplicationName)
		e.Str("server_address", pyroCfg.ServerAddress)
		e.Dur("upload_rate", pyroCfg.UploadRate)
		e.Strs("profile_types", profileTypeNames(pyroCfg.ProfileTypes))
		e.Any("tags", pyroCfg.Tags)
		e.Bool("disable_gc_runs", pyroCfg.DisableGCRuns)
		e.Bool("basic_auth", pyroCfg.BasicAuthUser != "")
		e.Bool("tenant_id", pyroCfg.TenantID != "")
		e.Bool("lifecycle_hook", p.Lc != nil)
	}).Msg("starting pyroscope profiler")

	profiler := assert.Must1(pyroscope.Start(pyroCfg))
	logger.Info().
		Str("application_name", pyroCfg.ApplicationName).
		Str("server_address", pyroCfg.ServerAddress).
		Msg("pyroscope profiler started")

	if p.Lc != nil {
		logger.Debug().Msg("register pyroscope profiler lifecycle stop hook")
		p.Lc.BeforeStop(lifecycle.WrapNoCtxErr(func() {
			stopProfiler(logger, profiler)
		}))
	}

	return &Client{profiler: profiler, logger: logger}
}

func resolveLogger(p Param) log.Logger {
	if p.Logger != nil {
		return p.Logger.WithName(Name)
	}
	return log.GetLogger(Name)
}

func stopProfiler(logger log.Logger, profiler *pyroscope.Profiler) {
	logger.Info().Msg("stopping pyroscope profiler")
	if err := profiler.Stop(); err != nil {
		logger.Err(err).Msg("failed to stop pyroscope profiler")
		return
	}
	logger.Info().Msg("pyroscope profiler stopped")
}

func profileTypeNames(types []pyroscope.ProfileType) []string {
	return lo.Map(types, func(item pyroscope.ProfileType, _ int) string {
		return string(item)
	})
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
