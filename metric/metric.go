package metric

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/lifecycle"
	"github.com/pubgo/funk/merge"
	"github.com/pubgo/funk/version"
	"github.com/uber-go/tally/v4"
)

func New(m lifecycle.Lifecycle, cfg *Cfg, optMap map[string]*tally.ScopeOptions) Metric {
	cfg = merge.Struct(generic.Ptr(DefaultCfg()), cfg).Unwrap()
	var opts = optMap[cfg.Driver]
	if opts == nil {
		return tally.NoopScope
	}

	opts.Tags = Tags{"project": version.Project()}
	if cfg.Separator != "" {
		opts.Separator = cfg.Separator
	}

	scope, closer := tally.NewRootScope(*opts, cfg.Interval)
	m.BeforeStop(func() { assert.Must(closer.Close()) })

	registerVars(scope)
	return scope
}
