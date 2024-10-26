package log

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog"
)

const (
	ModuleName = "module"
)

type (
	Map           = map[string]any
	Hook          = zerolog.Hook
	Level         = zerolog.Level
	EnableChecker = func(ctx context.Context, lvl Level, nameOrMessage string, fields Map) bool
)

type Logger interface {
	WithName(name string) Logger
	WithFields(m Map) Logger
	WithCallerSkip(skip int) Logger
	WithEvent(evt Event) Logger
	WithLevel(lvl Level) Logger

	Debug(ctx ...context.Context) Event
	Info(ctx ...context.Context) Event
	Warn(ctx ...context.Context) Event
	Error(ctx ...context.Context) Event
	Err(err error, ctx ...context.Context) Event
	Panic(ctx ...context.Context) Event
	Fatal(ctx ...context.Context) Event

	nameWithCaller(name string, caller int) Logger
}

type StdLogger interface {
	Printf(format string, v ...interface{})
	Logf(format string, v ...interface{})
	Print(v ...interface{})
	Log(v ...interface{})
	Println(v ...interface{})
}

type Event interface {
	Dict(key string, dict Event) Event
	Enabled() bool
	Discard() Event
	Msg(msg string)
	Msgf(format string, v ...any)
	Fields(fields interface{}) Event
	MACAddr(key string, ha net.HardwareAddr) Event
	IPPrefix(key string, pfx net.IPNet) Event
	IPAddr(key string, ip net.IP) Event
	Caller(skip ...int) Event
	CallerSkipFrame(skip int) Event
	Type(key string, val interface{}) Event
	Any(key string, i any) Event
	Durs(key string, d []time.Duration) Event
	Dur(key string, d time.Duration) Event
	Times(key string, t []time.Time) Event
	Time(key string, t time.Time) Event
	Timestamp() Event
	Floats64(key string, f []float64) Event
	Float64(key string, f float64) Event
	Floats32(key string, f []float32) Event
	Float32(key string, f float32) Event
	Uints64(key string, i []uint64) Event
	Uint64(key string, i uint64) Event
	Err(err error) Event
	Ctx(ctx context.Context) Event
	Func(func(e Event)) Event
	Str(key string, value string) Event
	Int(key string, value int) Event
}
