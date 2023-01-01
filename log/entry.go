package log

import (
	"github.com/pubgo/funk/log/log_fields"
	"net"
	"time"

	_ "github.com/phuslu/log"

	"github.com/pubgo/funk/logger"
	"github.com/pubgo/funk/stack"
)

func newEntry(log *loggerImpl, level logger.Level) logger.Entry {
	return &entryImpl{
		log:         log,
		level:       level,
		callerDepth: log.callerDepth,
	}
}

var _ logger.Entry = (*entryImpl)(nil)

type entryImpl struct {
	log         *loggerImpl
	callerDepth int
	fields      logger.Fields
	level       logger.Level
	msg         string
	time        time.Time
}

func (e *entryImpl) Base64(dst, val []byte) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Caller(depth ...int) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) AddXID() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) AddGoID() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) GetCaller() int {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) GetTime() time.Time {
	return e.time
}

func (e *entryImpl) GetMsg() string {
	return e.msg
}

func (e *entryImpl) WithCallerDepth(depth int) logger.Entry {
	e.callerDepth += depth
	return e
}

func (e *entryImpl) GetCallerDepth() int {
	return e.callerDepth
}

func (e *entryImpl) GetName() string {
	return e.log.name
}

func (e *entryImpl) GetLevel() logger.Level {
	return e.level
}

func (e *entryImpl) GetFields() []logger.Field {
	return e.fields
}

func (e *entryImpl) Log(msg string) {

}

func (e *entryImpl) Logf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Time(key string, val time.Time) logger.Entry {
	if e == nil {
		return nil
	}
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Fields(fields ...logger.Field) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int(key string, val int) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) IntL(key string, val ...int) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int8(key string, val int8) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int8L(key string, val ...int8) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int16(key string, val int16) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int16L(key string, val ...int16) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int32(key string, val int32) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int32L(key string, val ...int32) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int64(key string, val int64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Int64L(key string, val ...int64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint(key string, val uint) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) UintL(key string, val ...uint) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint8(key string, val uint8) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint8L(key string, val ...uint8) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint16(key string, val uint16) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint16L(key string, val ...uint16) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint32(key string, val uint32) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint32L(key string, val ...uint32) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint64(key string, val uint64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Uint64L(key string, val ...uint64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Float32(key string, val float32) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Float32L(key string, val ...float32) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Float64(key string, val float64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Float64L(key string, val ...float64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Any(key string, val interface{}) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Json(key string, val interface{}) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Hex(dst, val []byte) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Bool(key string, val bool) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) BoolL(key string, val ...bool) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Bytes(key string, val []byte) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) JsonRaw(key string, val []byte) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Dur(key string, dur time.Duration) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) DurL(key string, val ...time.Duration) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) TimeL(key string, val ...time.Time) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) IPAddr(key string, ip net.IP) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) IPAddrL(key string, ip ...net.IP) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) IPPrefix(key string, pfx net.IPNet) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) IPPrefixL(key string, pfx ...net.IPNet) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) MACAddr(key string, ha net.HardwareAddr) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) MACAddrL(key string, ha ...net.HardwareAddr) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Discard() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Print(msg string) {
	e.time = time.Now().UTC()
	e.msg = msg
	e.fields = append(e.fields, log_fields.Str("caller", stack.Caller(e.callerDepth).String()))

	var err = writer.WriteEntry(e)
}

func (e *entryImpl) Printf(format string, args ...any) {
	stack.Caller(e.callerDepth + 1)
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) TraceLog(msg string) {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) TraceLogf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Str(key string, value string) logger.Entry {
	e.fields = append(e.fields, &field{name: key, value: value, kind: log_types.String})
	return e
}

func (e *entryImpl) StrL(key string, value ...string) logger.Entry {
	e.fields = append(e.fields, &field{name: key, value: value, kind: log_types.StringL})
	return e
}

func (e *entryImpl) Field(field ...logger.Field) logger.Entry {
	e.fields = append(e.fields, field...)
	return e
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (e *entryImpl) Trace(msg string) *Entry {
	e.Info(msg)
	v := e.WithFields(e.Fields)
	v.Message = msg
	v.start = time.Now()
	return v
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the log level is error.
func (e *entryImpl) Stop(err *error) {
	if err == nil || *err == nil {
		e.WithDuration(time.Since(e.start)).Info(e.Message)
	} else {
		e.WithDuration(time.Since(e.start)).WithError(*err).Error(e.Message)
	}
}
