package errors

import (
	"reflect"
	"sync"

	"github.com/pubgo/funk/stack"
)

var skipStackMap sync.Map

func RegStackPkgFilter(fn ...interface{}) {
	if len(fn) == 0 {
		skipStackMap.Store(stack.Caller(1).Pkg, nil)
		return
	}

	for i := range fn {
		if reflect.TypeOf(fn[i]).Kind() == reflect.String {
			skipStackMap.Store(fn[i].(string), nil)
		} else {
			skipStackMap.Store(stack.CallerWithFunc(fn[0]).Pkg, nil)
		}
	}
}
