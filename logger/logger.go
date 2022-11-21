package logger

import "context"

type Tagger interface {
	Key() string
	Value() interface{}
}

type Logger interface {
	Info(level uint, msg string, tags ...Tagger)
	Error(err error, msg string, tags ...Tagger)
}

type Hook interface {
	Hook(tags []Tagger) []Tagger
}

type CtxParser func(ctx context.Context) (bool, []Tagger)
type ValueParser func(v interface{}) (bool, string)

type Marshaler interface {
	// MarshalLog can be used to:
	//   - ensure that structs are not logged as strings when the original
	//     value has a String method: return a different type without a
	//     String method
	//   - select which fields of a complex type should get logged:
	//     return a simpler struct with fewer fields
	//   - log unexported fields: return a different struct
	//     with exported fields
	//
	// It may return any value of any type.
	MarshalLog() interface{}
}
