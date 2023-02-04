package tracing

import (
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/vars"
)

func init() {
	vars.Register(Name+"_factory", func() interface{} {
		var data = make(map[string]string)
		factories.Range(func(key, value interface{}) bool {
			data[key.(string)] = stack.CallerWithFunc(value).String()
			return true
		})
		return data
	})
}
