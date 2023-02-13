package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/runmode"
	"github.com/pubgo/funk/tracing"
	"go.etcd.io/etcd/api/v3/version"
)

func init() {
	tracing.RegisterFactory(Name, func(cfgMap config.CfgMap) error {
		tracing.GetSpanID = GetSpanID

		var cfg = DefaultCfg()
		cfg.ServiceName = runmode.Project
		cfg.Tags = append(cfg.Tags, opentracing.Tag{Key: "version", Value: version.Version})
		assert.Must(cfgMap.Decode(&cfg))
		return New(cfg)
	})
}
