package anyhow

type Catchable interface {
	Catch(err *error) bool
	CatchErr(err *Error) bool
}

// Checkable defines types that can be checked for Ok/Error state
type Checkable interface {
	IsOk() bool
	IsErr() bool
	String() string
}
