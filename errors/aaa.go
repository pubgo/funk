package errors

import (
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"github.com/rs/zerolog"
)

type Event = zerolog.Event
type Tags map[string]any
type Errors []error

type ErrUnwrap interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type Error interface {
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
	BizCode() string
	Reason() string
	Code() errorpb.Code
	Tags() map[string]string
	AddTag(key string, val string)
	SetCode(code errorpb.Code)
	SetReason(reason string)
	SetBizCode(biz string)
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
