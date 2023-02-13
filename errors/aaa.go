package errors

import (
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"github.com/rs/zerolog"
)

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
	Name() string
	Reason() string
	Code() errorpb.Code
	Status() uint32
	Tags() map[string]string

	SetErr(err error) ErrCode
	AddTag(key string, val string) ErrCode
	SetCode(code errorpb.Code) ErrCode
	SetReason(reason string) ErrCode
	SetName(biz string) ErrCode
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
