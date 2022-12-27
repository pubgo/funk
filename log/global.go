package log

import (
	"context"
	"github.com/pubgo/funk/logger/logs"

	"github.com/pubgo/funk/logger"
)

type logCtx struct{}

var writer logger.Writer
var hooks []logger.Hook
var gv = logger.DEBUG
var log Logger

func init() {
	SetWriter(logs.NewTextLog())
}

func New(name string) Logger {
	return Named(name)
}

func AddHook(h logger.Hook) {
	if h == nil {
		panic("log: log hook is nil")
	}

	hooks = append(hooks, h)
}

func SetWriter(w logger.Writer) {
	if w == nil {
		panic("log: log writer is nil")
	}

	writer = w
}

func SetLevel(level logger.Level) {
	gv = level
}

func GetLevel() logger.Level {
	return gv
}

func Ctx(ctx context.Context) Logger {
	if ctx != nil {
		if l, ok := ctx.Value(logCtx{}).(Logger); ok {
			return l
		}
	}
	return log
}

func Named(name string) Logger {
	return log.WithName(name)
}

// Trace starts a new message with trace level.
func Trace() (e *Entry) {
	if DefaultLogger.silent(TraceLevel) {
		return nil
	}
	e = DefaultLogger.header(TraceLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Debug starts a new message with debug level.
func Debug() (e *Entry) {
	if DefaultLogger.silent(DebugLevel) {
		return nil
	}
	e = DefaultLogger.header(DebugLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Info starts a new message with info level.
func Info() (e *Entry) {
	if DefaultLogger.silent(InfoLevel) {
		return nil
	}
	e = DefaultLogger.header(InfoLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Warn starts a new message with warning level.
func Warn() (e *Entry) {
	if DefaultLogger.silent(WarnLevel) {
		return nil
	}
	e = DefaultLogger.header(WarnLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Error starts a new message with error level.
func Error() (e *Entry) {
	if DefaultLogger.silent(ErrorLevel) {
		return nil
	}
	e = DefaultLogger.header(ErrorLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Fatal starts a new message with fatal level.
func Fatal() (e *Entry) {
	if DefaultLogger.silent(FatalLevel) {
		return nil
	}
	e = DefaultLogger.header(FatalLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Panic starts a new message with panic level.
func Panic() (e *Entry) {
	if DefaultLogger.silent(PanicLevel) {
		return nil
	}
	e = DefaultLogger.header(PanicLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	return
}

// Printf sends a log entry without extra field. Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	e := DefaultLogger.header(noLevel)
	if caller, full := DefaultLogger.Caller, false; caller != 0 {
		if caller < 0 {
			caller, full = -caller, true
		}
		var rpc [1]uintptr
		e.caller(callers(caller, rpc[:]), rpc[:], full)
	}
	e.Msgf(format, v...)
}
