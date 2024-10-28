package log

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/v2/log/log_internal"
)

// GetLogger get global log
func GetLogger(name string) Logger {
	return log_internal.GetLogger().WithNameCaller(name, 1)
}

func Err(err error, ctx ...context.Context) Event {
	return log_internal.GetLogger().Err(err, ctx...)
}

func Debug(ctx ...context.Context) Event {
	return log_internal.GetLogger().Debug(ctx...)
}

func Info(ctx ...context.Context) Event {
	return log_internal.GetLogger().Info(ctx...)
}

func Warn(ctx ...context.Context) Event {
	return log_internal.GetLogger().Warn(ctx...)
}

func Error(ctx ...context.Context) Event {
	return log_internal.GetLogger().Error(ctx...)
}

func Fatal(ctx ...context.Context) Event {
	return log_internal.GetLogger().Fatal(ctx...)
}

func Panic(ctx ...context.Context) Event {
	return log_internal.GetLogger().Panic(ctx...)
}

func Print(v ...interface{}) {
	log_internal.GetLogger().Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	log_internal.GetLogger().Debug().CallerSkipFrame(1).Msgf(format, v...)
}
