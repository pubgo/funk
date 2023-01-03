package errors

import (
	"fmt"
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

type Map map[string]any

type TagWrapper interface {
	Error
	Tags() map[string]any
}

type MsgWrapper interface {
	Error
	Msg() string
}

type CodeWrapper interface {
	Error
	Code() codes.Code
}

type CallerWrapper interface {
	Error
	Caller() *stack.Frame
}

type StackWrapper interface {
	Error
	Stack() []*stack.Frame
}

type BizCodeWrapper interface {
	Error
	BizCode() string
}

type Error interface {
	Error() string
	String() string
	Format(s fmt.State, verb rune)
	Unwrap() error
	As(target interface{}) bool
	MarshalJSON() ([]byte, error)
}
