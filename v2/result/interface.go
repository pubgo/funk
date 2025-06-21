package result

import "context"

type Catchable interface {
	Catch(err *error, contexts ...context.Context) bool
	CatchErr(err ErrSetter, contexts ...context.Context) bool
}

// Checkable defines types that can be checked for Ok/Error state
type Checkable interface {
	IsOK() bool
	IsErr() bool
	GetErr() error
	String() string
}

type ErrSetter interface {
	Checkable
	setError(err error)
}

type UnWrapper[T any] interface {
	Unwrap(setter *error, contexts ...context.Context) T
	UnwrapErr(setter ErrSetter, contexts ...context.Context) T
}
