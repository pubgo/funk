package errors

import (
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

type Map map[string]any
type XErr interface {
	_err()
	Error() string
	String() string
	Code() codes.Code
	BizCode() string
	Unwrap() error
	As(target interface{}) bool
	MarshalJSON() ([]byte, error)
	AddCaller() XErr
	AddTag(key string, value any) XErr
	AddCode(code codes.Code) XErr
	AddBizCode(biz string) XErr
	AddMsg(msg string) XErr
	WithTag(key string, value any) XErr
	WithCode(code codes.Code) XErr
	WithBizCode(biz string) XErr
}

type Status interface {
	Code() codes.Code
}

type StackTrace interface {
	Stack() *stack.Frame
}

type RespErr interface {
	XErr
}
