package log

import (
	"context"
	_ "github.com/phuslu/log"
	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog"
)

type Map = map[string]any
type Hook = zerolog.Hook
type Logger interface {
	WithName(name string) Logger
	WithFields(m Map) Logger
	WithHooks(hooks ...Hook) Logger
	WithCallerSkip(skip int) Logger
	WithCtx(ctx context.Context) context.Context
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Err(err error) *zerolog.Event
	Panic() *zerolog.Event
	Fatal() *zerolog.Event
	Print(msg string)
	Printf(format string, args ...any)
}
