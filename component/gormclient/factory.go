package gormclient

import (
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/config"
	"github.com/pubgo/funk/recovery"
	"gorm.io/gorm"
)

type Factory func(cfg config.Node) gorm.Dialector

var factories = make(map[string]Factory)

func Get(name string) Factory  { return factories[name] }
func List() map[string]Factory { return factories }
func Register(name string, broker Factory) {
	defer recovery.Exit()
	assert.If(name == "" || broker == nil, "[broker,name] should not be null")
	assert.If(factories[name] != nil, "[broker] %s already exists", name)
	factories[name] = broker
}
