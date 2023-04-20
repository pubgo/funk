package errors

import (
	"fmt"

	"github.com/pubgo/funk/stack"
)

var _ Error = (*ErrWrap)(nil)
var _ fmt.Formatter = (*ErrWrap)(nil)

type ErrWrap struct {
	err    error
	caller *stack.Frame
	stack  []*stack.Frame
	fields map[string]any
}

func (e ErrWrap) Format(f fmt.State, verb rune) {
	//TODO implement me
	panic("implement me")
}

func (e ErrWrap) Kind() string {
	return "err_wrap"
}

func (e ErrWrap) Error() string {
	//TODO implement me
	panic("implement me")
}

func (e ErrWrap) String() string {
	//TODO implement me
	panic("implement me")
}

func (e ErrWrap) Unwrap() error {
	//TODO implement me
	panic("implement me")
}

func (e ErrWrap) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
