package errors

import (
	"reflect"
	"sync"

	"github.com/pubgo/funk/stack"
)

var skipStackMap sync.Map

// RegStackPkgFilter filter fn , pkg
func RegStackPkgFilter(filter ...interface{}) {
	if len(filter) == 0 {
		skipStackMap.Store(stack.Caller(1).Pkg, nil)
		return
	}

	for _, ff := range filter {
		if ff == nil {
			continue
		}

		switch ff.(type) {
		case reflect.Type:
			skipStackMap.Store(stack.CallerWithType(ff.(reflect.Type)).Pkg, nil)
			continue
		}

		typ := reflect.TypeOf(ff)
		switch typ.Kind() {
		case reflect.String:
			skipStackMap.Store(ff.(string), nil)
		case reflect.Func:
			skipStackMap.Store(stack.CallerWithFunc(ff).Pkg, nil)
		default:
		}
	}
}
