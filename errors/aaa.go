package errors

type Tags map[string]any

func (t Tags) Set(key string, val any) {
	t[key] = val
}

type Tag struct {
	k string
	v any
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

type Errors interface {
	Error
	Errors() []error
	Append(err error) error
}

type Error interface {
	Kind() string
	Error() string
	String() string
	MarshalJSON() ([]byte, error)
}
