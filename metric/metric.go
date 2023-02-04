package metric

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/lifecycle"
	"github.com/pubgo/funk/runmode"
	"github.com/uber-go/tally/v4"
)

func New(m lifecycle.Lifecycle, cfg *Cfg, optMap map[string]*tally.ScopeOptions) Metric {
	var opts = optMap[cfg.Driver]
	if opts == nil {
		opts = &tally.ScopeOptions{Reporter: tally.NullStatsReporter}
	}

	opts.Tags = Tags{"service": runmode.Project}
	if cfg.Separator != "" {
		opts.Separator = cfg.Separator
	}

	scope, closer := tally.NewRootScope(*opts, cfg.Interval)
	m.BeforeStop(func() { assert.Must(closer.Close()) })

	registerVars(scope)
	return scope
}
