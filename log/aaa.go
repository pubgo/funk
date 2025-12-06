package log

import (
	"context"

	"github.com/rs/zerolog"
)

type (
	Fields        = map[string]any
	Hook          = zerolog.Hook
	Event         = zerolog.Event
	Level         = zerolog.Level
	EnableChecker = func(ctx context.Context, lvl Level, name, message string, fields Fields) bool
)

type Logger interface {
	WithName(name string) Logger
	WithFields(m Fields) Logger
	WithCallerSkip(skip int) Logger
	WithLevel(lvl Level) Logger

	Debug(ctx ...context.Context) *Event
	Info(ctx ...context.Context) *Event
	Warn(ctx ...context.Context) *Event
	Error(ctx ...context.Context) *Event
	Err(err error, ctx ...context.Context) *Event
	Panic(ctx ...context.Context) *Event
	Fatal(ctx ...context.Context) *Event

	nameWithCaller(name string, caller int) Logger
}

type StdLogger interface {
	Printf(format string, v ...any)
	Logf(format string, v ...any)
	Print(v ...any)
	Log(v ...any)
	Println(v ...any)
}
