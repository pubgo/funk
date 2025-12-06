package errors

import (
	"fmt"
)

type ErrorID interface {
	error
	ID() string
}

type Error interface {
	ErrorID
	String() string
	MarshalJSON() ([]byte, error)
}

type ErrUnwrapper interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type Tags map[string]any

func (t Tags) ToMapString() map[string]string {
	data := make(map[string]string, len(t))
	for key, value := range t {
		data[key] = fmt.Sprintf("%v", value)
	}
	return data
}
