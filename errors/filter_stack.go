package errors

import (
	"sync"

	"github.com/pubgo/funk/stack"
)

var skipStackMap sync.Map

func RegStackPkgFilter(fn ...interface{}) {
	var s *stack.Frame
	if len(fn) == 0 || fn[0] == nil {
		s = stack.Caller(1)
	} else {
		s = stack.CallerWithFunc(fn[0])
	}
	skipStackMap.Store(s.Pkg, nil)
}
