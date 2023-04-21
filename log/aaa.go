package log

import (
	"github.com/rs/zerolog"
)

type Map = map[string]any
type Hook = zerolog.Hook
type Event = zerolog.Event
type Level = zerolog.Level

type Logger interface {
	WithName(name string) Logger
	WithFields(m Map) Logger
	WithHooks(hooks ...Hook) Logger
	WithCallerSkip(skip int) Logger
	WithEvent(evt *Event) Logger
	WithLevel(lvl Level) Logger
	Debug() *Event
	Info() *Event
	Warn() *Event
	Error() *Event
	Err(err error) *Event
	Panic() *Event
	Fatal() *Event
}

type StdLogger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
}
