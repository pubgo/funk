package errors

import (
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

type Map map[string]any

type Errors []error

type errUnwrap interface {
	Unwrap() error
}

type opaqueWrapper struct {
}

type errIs interface {
	Is(error) bool
}

type errAs interface {
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

type XError interface {
	Error
	BizCode() string
	Stack() []*stack.Frame
	Code() codes.Code
	Msg() string
	Tags() map[string]any
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

// T is a shortcut to make a Tag
func T(key string, value any) Tag {
	return Tag{Key: key, Value: value}
}

// Tag contains a single key value combination
// to be attached to your error
type Tag struct {
	Key   string
	Value any
}

//func RegisterHelper(helper Helper) {
//	for i := 0; i < len(helpers); i++ {
//		if reflect.ValueOf(helpers[i]).Pointer() == reflect.ValueOf(helper).Pointer() {
//			return
//		}
//	}
//	helpers = append(helpers, helper)
//}
//
//type Helper func(Chain, error) bool
