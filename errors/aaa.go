package errors

import (
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

type Map map[string]any

type ErrUnwrap interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type ITagWrap interface {
	Error
	Tags() map[string]any
}

type IMsgWrap interface {
	Error
	Msg() string
}

type ICodeWrap interface {
	Error
	Code() codes.Code
}

type IStackWrap interface {
	Error
	Stack() []*stack.Frame
}

type IBizCodeWrap interface {
	Error
	BizCode() string
}

type Error interface {
	Error() string
	String() string
	Unwrap() error
	MarshalJSON() ([]byte, error)
}

type RespErr struct {
	Cause   error          `json:"cause"`
	Msg     string         `json:"msg"`
	Code    codes.Code     `json:"code"`
	BizCode string         `json:"biz_code"`
	Tags    map[string]any `json:"tags"`
}
