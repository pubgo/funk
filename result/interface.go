package result

import (
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

type Void = funk.Void
