package log

import (
	"context"

	"github.com/pubgo/funk/logger"
)

type logCtx struct{}

var writer logger.Logger
var hooks []logger.Hook
var ll uint
var log Logger
var NameDelim = "."

func AddHook(h logger.Hook) {
	if h == nil {
		panic("log: log hook is nil")
	}

	hooks = append(hooks, h)
}

func SetWriter(w logger.Logger) {
	if w == nil {
		panic("log: log writer is nil")
	}

	writer = w
}

func SetLevel(level uint) {
	ll = level
}

func GetLevel() uint {
	return ll
}

func Ctx(ctx context.Context) Logger {
	if ctx != nil {
		if l, ok := ctx.Value(logCtx{}).(Logger); ok {
			return l
		}
	}
	return log
}

func V(level uint) Logger {
	return log.V(level)
}

func Named(name string) Logger {
	return log.WithName(name)
}
