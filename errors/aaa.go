package errors

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

type Errors interface {
	Error
	Errors() []error
	Append(err error) error
}

type Error interface {
	Kind() string
	Error() string
	String() string
	Unwrap() error
	MarshalJSON() ([]byte, error)
}
