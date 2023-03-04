package drivers

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/metric"
	"github.com/pubgo/funk/recovery"
	"github.com/uber-go/tally/v4"
)

type Factory func(cfg *metric.Cfg, log log.Logger) *tally.ScopeOptions

var factories = make(map[string]Factory)

func Get(name string) Factory  { return factories[name] }
func List() map[string]Factory { return factories }
func Register(name string, broker Factory) {
	defer recovery.Exit()
	assert.If(name == "" || broker == nil, "[broker,name] should not be null")
	assert.If(factories[name] != nil, "[broker] %s already exists", name)
	factories[name] = broker
}
