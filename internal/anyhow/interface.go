package anyhow

import "context"

type Catchable interface {
	Catch(err *error, contexts ...context.Context) bool
	CatchErr(err *Error, contexts ...context.Context) bool
}

// Checkable defines types that can be checked for Ok/Error state
type Checkable interface {
	IsOK() bool
	IsErr() bool
	GetErr() error
	String() string
}
