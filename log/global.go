package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/generic"
)

var (
	logEnableChecker = func(ctx context.Context, lvl Level, nameOrMessage string, fields Map) bool { return true }
	zErrMarshalFunc  = zerolog.ErrorMarshalFunc
	logGlobalHook    = zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		if logEnableChecker == nil {
			return
		}

		ctx := e.GetCtx()
		if logEnableChecker(ctx, level, message, getFieldFromCtx(ctx)) {
			return
		}

		e.Discard()
	})
	_ = generic.Init(func() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.ErrorMarshalFunc = func(err error) interface{} {
			if err == nil {
				return nil
			}

			switch err.(type) {
			case json.Marshaler:
				return &logLogObjectMarshaler{err: err}
			}

			if zErrMarshalFunc == nil {
				return fmt.Sprintf("%v", err)
			}

			return zErrMarshalFunc(err)
		}
	})

	// stdZeroLog default zerolog for debug
	stdZeroLog = generic.Ptr(
		zerolog.New(os.Stderr).
			Level(zerolog.DebugLevel).
			With().Timestamp().
			Caller().Logger().
			Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
				w.Out = os.Stderr
				w.TimeFormat = time.RFC3339
			})).Hook(new(hookImpl), logGlobalHook),
	)

	_ = generic.Init(func() {
		zlog.Logger = generic.FromPtr(stdZeroLog)
	})

	// stdLog is the global logger.
	stdLog = New(nil)
)

// GetLogger get global log
func GetLogger(names ...string) Logger {
	if len(names) == 0 || names[0] == "" {
		return stdLog
	}
	return stdLog.nameWithCaller(names[0], 1)
}

// SetLogger set global log
func SetLogger(log *zerolog.Logger) {
	assert.If(log == nil, "[log] should not be nil")

	log = generic.Ptr(log.Hook(logGlobalHook))

	stdZeroLog = log
	zlog.Logger = *log
}

func SetEnableChecker(checker EnableChecker) {
	if checker == nil {
		return
	}

	logEnableChecker = checker
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

func Output(w io.Writer) Logger {
	return New(generic.Ptr(stdZeroLog.Output(w)))
}

type writerFunc func(p []byte) (n int, err error)

func (w writerFunc) Write(p []byte) (n int, err error) {
	return w(p)
}

func OutputWriter(w func(p []byte) (n int, err error)) Logger {
	return New(generic.Ptr(stdZeroLog.Output(writerFunc(w))))
}
