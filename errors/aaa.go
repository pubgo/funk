package errors

import (
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

type Tags map[string]any

type Errors []error

type errUnwrap interface {
	Unwrap() error
}

type opaqueWrapper struct {
	err error
}

type errIs interface {
	Is(error) bool
}

type errAs interface {
	As(any) bool
}

type XErr = XError
type XError interface {
	Error() string
	String() string
	Unwrap() error
	MarshalJSON() ([]byte, error)

	BizCode() string
	Stack() []*stack.Frame
	Code() codes.Code
	Msg() string
	Tags() Tags

	AddBizCode(biz string)
	AddStack()
	AddMsg(msg string)
	AddCode(code codes.Code)
	AddTag(key string, val any)
	AddTags(m Tags)
}

type RespErr struct {
	Cause   error          `json:"cause"`
	Msg     string         `json:"msg"`
	Code    codes.Code     `json:"code"`
	BizCode string         `json:"biz_code"`
	Tags    map[string]any `json:"tags"`
}
