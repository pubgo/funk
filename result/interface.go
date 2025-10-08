package result

import (
	"context"

	"github.com/pubgo/funk/v2"
)

// Checkable defines types that can be checked for Ok/Error state
type Checkable interface {
	IsOK() bool
	IsErr() bool
	GetErr() error
	String() string
}

type ErrSetter interface {
	Checkable
	setErrorInner()
}

type Catchable interface {
	CatchErr(err *error, contexts ...context.Context) bool
	Catch(err ErrSetter, contexts ...context.Context) bool
}

type UnWrapper[T any] interface {
	UnwrapErr(setter *error, contexts ...context.Context) T
	Unwrap(setter ErrSetter, contexts ...context.Context) T
}

type Void = funk.Void
