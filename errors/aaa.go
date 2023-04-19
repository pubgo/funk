package errors

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

// base error
// with message
// with metadata
// with tags
// 一个底层error，用于展示，查询和和打印log
// 一个wrap error，wrap一些主要的信息，用于参考和追踪错误链路
// 一个code error
// 错误是要能够被打印和输出的
// 错误是要能够被识别和处理的
func NewError(prefix string) {

}

type ErrCode1 struct {
	code *errorpb.ErrCode
}

type Event = zerolog.Event
type Tags map[string]any

type Errors interface {
	Error
	Errors() []error
	Append(err error) error
}

type ErrUnwrap interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type Error2 interface {
	error
	fmt.Stringer
	json.Marshaler
	WithField(key, value string) Error2
	SetTag(key, value string) Error2
	Event() *Event
}

type Error1 interface {
	error
	fmt.Stringer
	ErrUnwrap
	json.Marshaler
}

type Error interface {
	Kind() string
	Error() string
	String() string
	Unwrap() error
	MarshalJSON() ([]byte, error)
}

type ErrEvent interface {
	Error
	Event() *Event
	AddEvent(evt *Event)
}

type ErrCode interface {
	Error
	Reason() string
	Code() errorpb.Code
	Status() string
	Tags() map[string]string

	SetErr(err error) ErrCode
	AddTag(key string, val string) ErrCode
	SetCode(code errorpb.Code) ErrCode
	SetStatus(status string) ErrCode
	SetReason(reason string) ErrCode
}

type ErrStack interface {
	Error
	AddStack()
	Stack() []*stack.Frame
}

// event 和<zerolog.Event>内存对齐
type event struct {
	buf []byte
}
