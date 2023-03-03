package prometheus

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/debug"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/metric"
	tally "github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
)

const Name = "prometheus"
const urlPath = "/metrics"

func New(conf *metric.Cfg, log log.Logger) map[string]*tally.ScopeOptions {
	if conf.Driver != Name {
		return nil
	}

	opts := tally.ScopeOptions{}
	opts.Separator = prometheus.DefaultSeparator
	opts.SanitizeOptions = &prometheus.DefaultSanitizerOpts

	var proCfg = &prometheus.Configuration{TimerType: "histogram"}

	if conf.DriverCfg != nil {
		assert.Must(conf.DriverCfg.Decode(proCfg))
	}

	var logs = log.WithName(metric.Name).WithName(Name)
	reporter := assert.Must1(proCfg.NewReporter(
		prometheus.ConfigurationOptions{
			OnError: func(err error) {
				logs.Err(err).Any("metric-config", conf).Msg("metric.prometheus init error")
			},
		},
	))
	debug.Get(urlPath, debug.Wrap(reporter.HTTPHandler()))

	opts.CachedReporter = reporter
	return map[string]*tally.ScopeOptions{Name: &opts}
}
