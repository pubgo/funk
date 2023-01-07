package log

import (
	"context"

	"github.com/rs/zerolog"
)

type Map = map[string]any
type Hook = zerolog.Hook
type Event = zerolog.Event
type Logger interface {
	WithName(name string) Logger
	WithFields(m Map) Logger
	WithHooks(hooks ...Hook) Logger
	WithCallerSkip(skip int) Logger
	WithCtx(ctx context.Context) context.Context
	Debug() *Event
	Info() *Event
	Warn() *Event
	Error() *Event
	Err(err error) *Event
	Panic() *Event
	Fatal() *Event
	Print(msg string)
	Printf(format string, args ...any)
}
