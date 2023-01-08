package errors

import (
	"github.com/pubgo/funk/stack"
)

func (t *baseErr) Stack() []*stack.Frame {
	return t.stacks
}

func (t *baseErr) AddStack() {
	for i := 0; ; i++ {
		var cc = stack.Caller(1 + i)
		if cc == nil {
			break
		}

		if cc.IsRuntime() {
			continue
		}

		t.stacks = append(t.stacks, cc)
	}
}
