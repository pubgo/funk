package errors

import (
	"github.com/pubgo/funk/stack"
)

func getStack() (ss []*stack.Frame) {
	for i := 0; ; i++ {
		var cc = stack.Caller(1 + i)
		if cc == nil {
			break
		}

		if cc.IsRuntime() {
			continue
		}

		if _, ok := skipStack.Load(cc.Pkg); ok {
			continue
		}

		ss = append(ss, cc)
	}
	return ss
}
