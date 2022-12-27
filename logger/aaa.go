package logger

import (
	"context"
	"runtime"
	"time"
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

type Valuer func() string

type Tagger interface {
	Key() string
	Value() Valuer
}

type Writer interface {
	Close() error
	WriteEntry(level Level, topic string, caller int, tags []Tagger) error
}

type ObjectMarshaler interface {
	MarshalObject(e *Entry)
}

// Entry represents a log entry. It is instanced by one of the level method of Logger and finalized by the Msg or Msgf method.
type Entry struct {
	w *Writer

	buf    []byte
	topic  string
	caller int

	// Time for record log
	Time time.Time
	// Level log level for record
	Level Level
	// level name cache from Level
	levelName string

	// Fields custom fields data.
	// Contains all the fields set by the user.
	Fields Map
	// Data log context data
	Data Map
	// Extra log extra data
	Extra Map

	// Caller information
	Caller *runtime.Frame
	// CallerFlag value. default is equals to Logger.CallerFlag
	CallerFlag uint8
	// CallerSkip value. default is equals to Logger.CallerSkip
	CallerSkip int

	// Ctx context.Context
	Ctx context.Context
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
