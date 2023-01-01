package log

import (
	"context"
	"github.com/mattn/go-colorable"
	"golang.org/x/term"
	"os"

	"github.com/pubgo/funk/logger"
	"github.com/pubgo/funk/logger/logs"
)

type logCtx struct{}

var writer logger.Writer
var hooks []logger.Hook
var gv = logger.DEBUG
var stdLog *loggerImpl
var logNo uint64

func init() {
	SetWriter(logs.NewTextLog())
}

func AddHook(h logger.Hook) {
	if h == nil {
		panic("slog: log hook is nil")
	}

	hooks = append(hooks, h)
}

func SetWriter(w logger.Writer) {
	if w == nil {
		panic("slog: log writer is nil")
	}

	writer = w
}

func SetLevel(level logger.Level) {
	gv = level
}

func GetLevel() logger.Level {
	return gv
}

func Ctx(ctx context.Context) logger.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(logCtx{}).(logger.Logger); ok {
			return l
		}
	}
	return stdLog
}

func GetLogger(name string) logger.Logger {
	return stdLog.WithName(name)
}

// Trace starts a new message with trace level.
func Trace() (e logger.Entry) {
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
func Debug() (e logger.Entry) {
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
func Info() (e logger.Entry) {
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
func Warn() (e logger.Entry) {
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
func Error() (e logger.Entry) {
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
func Fatal() (e logger.Entry) {
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
func Panic() (e logger.Entry) {
	if !stdLog.Enabled(logger.CRITICAL) {
		return
	}

	stdLog.Panic()
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
	if !stdLog.Enabled(logger.INFO) {
		return
	}

	stdLog.Info().Caller(2).Printf(format, v...)
}

func init() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		StdoutHandler = StreamHandler(colorable.NewColorableStdout(), TerminalFormat())
	}

	if term.IsTerminal(int(os.Stderr.Fd())) {
		StderrHandler = StreamHandler(colorable.NewColorableStderr(), TerminalFormat())
	}
}
