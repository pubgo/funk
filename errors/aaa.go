package errors

import (
	"github.com/rs/zerolog"
)

type Event = zerolog.Event
type Tags map[string]any

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
	Kind() string
	Error() string
	String() string
	Unwrap() error
	MarshalJSON() ([]byte, error)
}

// event 和<zerolog.Event>内存对齐
type event struct {
	buf []byte
}
