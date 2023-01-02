package logger

import (
	"context"
	"net"
	"time"

	_ "github.com/phuslu/log"
	_ "github.com/rs/zerolog"
	_ "go.uber.org/zap"
)

type Logger interface {
	Enabled(level Level) bool
	WithName(name string) Logger
	WithFields(fields ...Field) Logger
	WithHooks(hooks ...Hook) Logger
	WithCaller(depth int) Logger
	WithCtx(ctx context.Context) context.Context
	Info() Entry
	Warn() Entry
	Error() Entry
	Err(err error) Entry
	Panic() Entry
	Fatal() Entry
	Print(msg string)
	Printf(format string, args ...any)
}

type Writer interface {
	Close() error
	Sync() error
	WriteEntry(ent Entry)
}

type Hook interface {
	Do(ent Entry) Entry
}

type Entry interface {
	Int(key string, val int) Entry
	IntL(key string, val ...int) Entry
	Int8(key string, val int8) Entry
	Int8L(key string, val ...int8) Entry
	Int16(key string, val int16) Entry
	Int16L(key string, val ...int16) Entry
	Int32(key string, val int32) Entry
	Int32L(key string, val ...int32) Entry
	Int64(key string, val int64) Entry
	Int64L(key string, val ...int64) Entry

	Uint(key string, val uint) Entry
	UintL(key string, val ...uint) Entry
	Uint8(key string, val uint8) Entry
	Uint8L(key string, val ...uint8) Entry
	Uint16(key string, val uint16) Entry
	Uint16L(key string, val ...uint16) Entry
	Uint32(key string, val uint32) Entry
	Uint32L(key string, val ...uint32) Entry
	Uint64(key string, val uint64) Entry
	Uint64L(key string, val ...uint64) Entry

	Float32(key string, val float32) Entry
	Float32L(key string, val ...float32) Entry
	Float64(key string, val float64) Entry
	Float64L(key string, val ...float64) Entry

	Json(key string, val interface{}) Entry
	JsonRaw(key string, val []byte) Entry
	Hex(dst, val []byte) Entry
	Base64(dst, val []byte) Entry
	Str(key string, val string) Entry
	StrL(key string, val ...string) Entry
	Bool(key string, val bool) Entry
	BoolL(key string, val ...bool) Entry
	Bytes(key string, val []byte) Entry
	Dur(key string, dur time.Duration) Entry
	DurL(key string, val ...time.Duration) Entry
	Time(key string, val time.Time) Entry
	TimeL(key string, val ...time.Time) Entry

	IPAddr(key string, ip net.IP) Entry
	IPAddrL(key string, ip ...net.IP) Entry
	IPPrefix(key string, pfx net.IPNet) Entry
	IPPrefixL(key string, pfx ...net.IPNet) Entry
	MACAddr(key string, ha net.HardwareAddr) Entry
	MACAddrL(key string, ha ...net.HardwareAddr) Entry

	BytesValue(key string, val Valuer) Entry
	BytesLValue(key string, val Valuer) Entry
	RawValue(key string, val Valuer) Entry
	RawLValue(key string, val Valuer) Entry
	Fields(fields ...Field) Entry

	WithCaller() Entry
	WithStack() Entry
	WithXID() Entry
	WithGoID() Entry
	WithCtx(ctx context.Context) Entry

	GetName() string
	GetLevel() Level
	GetFields() map[string]Field
	GetCtx() context.Context

	Log(msg string)
	Logf(format string, args ...interface{})
	Print(msg string)
	Printf(format string, args ...any)
	TraceLog(msg string)
	TraceLogf(format string, args ...any)
}
