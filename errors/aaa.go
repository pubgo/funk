package errors

import (
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

type Tags map[string]any

type Errors []error

type errUnwrap interface {
	Unwrap() error
}

type errIs interface {
	Is(error) bool
}

type errAs interface {
	As(any) bool
}

type XErr = XError
type Error interface {
	Error() string
	String() string
	Unwrap() error
	MarshalJSON() ([]byte, error)
}

type XError interface {
	Error
	BizCode() string
	Stack() []*stack.Frame
	Code() errorpb.Code
	Msg() string
	Tags() Tags

	AddBizCode(biz string)
	AddStack()
	AddMsg(msg string)
	AddCode(code errorpb.Code)
	AddTag(key string, val any)
	AddTags(m Tags)
}

type RespErr struct {
	Cause   error          `json:"cause"`
	Msg     string         `json:"msg"`
	Code    errorpb.Code   `json:"code"`
	BizCode string         `json:"biz_code"`
	Tags    map[string]any `json:"tags"`
}
