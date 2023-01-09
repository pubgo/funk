package log

import (
	"fmt"
	"os"
	"time"

	"github.com/pubgo/funk/generic"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

var (
	// stdZeroLog default zerolog just fro debug
	stdZeroLog = generic.Ptr(
		zerolog.New(os.Stderr).Level(zerolog.DebugLevel).
			With().Timestamp().Caller().Logger().
			Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
				w.Out = os.Stderr
				w.TimeFormat = time.RFC3339
			})),
	)

	// stdLog is the global logger.
	stdLog = New(nil).WithHooks(new(hookImpl))
)

func GetLogger(name string) Logger {
	return stdLog.WithName(name)
}

func SetLogger(log *zerolog.Logger) {
	if log == nil {
		panic("[log] should not be nil")
	}

	stdZeroLog = log
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return stdLog.Err(err)
}

// Debug starts a new message with debug level.
//
// You must call msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return stdLog.Debug()
}

// Info starts a new message with info level.
//
// You must call msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return stdLog.Info()
}

// Warn starts a new message with warn level.
//
// You must call msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return stdLog.Warn()
}

// Error starts a new message with error level.
//
// You must call msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return stdLog.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the msg method.
//
// You must call msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return stdLog.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return stdLog.Panic()
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
