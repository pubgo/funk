package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/generic"
)

var (
	_ = func() interface{} {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return nil
	}()

	// stdZerolog default zerolog for debug
	stdZerolog = generic.Ptr(
		zerolog.New(os.Stderr).Level(zerolog.DebugLevel).
			With().Timestamp().Caller().Logger().
			Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
				w.Out = os.Stderr
				w.TimeFormat = time.RFC3339
			})).Hook(new(hookImpl)),
	)

	// stdLog is the global logger.
	stdLog = New(nil)
)

// GetLogger get global log
func GetLogger(name string) Logger {
	return stdLog.WithName(name)
}

// SetLogger set global log
func SetLogger(log *zerolog.Logger) {
	assert.If(log == nil, "[log] should not be nil")
	stdZerolog = log
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call msg on the returned event in order to send the event.
func Err(err error, ctx ...context.Context) *zerolog.Event {
	return stdLog.Err(err, ctx...)
}

// Debug starts a new message with debug level.
//
// You must call msg on the returned event in order to send the event.
func Debug(ctx ...context.Context) *zerolog.Event {
	return stdLog.Debug(ctx...)
}

// Info starts a new message with info level.
//
// You must call msg on the returned event in order to send the event.
func Info(ctx ...context.Context) *zerolog.Event {
	return stdLog.Info(ctx...)
}

// Warn starts a new message with warn level.
//
// You must call msg on the returned event in order to send the event.
func Warn(ctx ...context.Context) *zerolog.Event {
	return stdLog.Warn(ctx...)
}

// Error starts a new message with error level.
//
// You must call msg on the returned event in order to send the event.
func Error(ctx ...context.Context) *zerolog.Event {
	return stdLog.Error(ctx...)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the msg method.
//
// You must call msg on the returned event in order to send the event.
func Fatal(ctx ...context.Context) *zerolog.Event {
	return stdLog.Fatal(ctx...)
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call msg on the returned event in order to send the event.
func Panic(ctx ...context.Context) *zerolog.Event {
	return stdLog.Panic(ctx...)
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	stdLog.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	stdLog.Debug().CallerSkipFrame(1).Msgf(format, v...)
}
