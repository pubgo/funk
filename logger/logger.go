package logger

import (
	"context"
)

type StdLogger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
}

type ExternalLogger interface {
	Print(...interface{})
	Println(...interface{})
	Error(...interface{})
	Warn(...interface{})
	Info(...interface{})
	Debug(...interface{})
}

type Tagger interface {
	Key() string
	Value() interface{}
}

type Logger interface {
	Log(tags ...Tagger)
}

type Hook interface {
	Hook(tags []Tagger) []Tagger
}

type CtxParser func(ctx context.Context) (bool, []Tagger)
type ValueParser func(v interface{}) (bool, string)
type Fields map[string]interface{}
type A any
type S []any
type Map []Field
type Field struct {
	Name  string
	Value interface{}
}

// F is a convenience constructor for Field.
func F(name string, value interface{}) Field {
	return Field{Name: name, Value: value}
}

// M is a convenience constructor for Map
func M(fs ...Field) Map {
	return fs
}
