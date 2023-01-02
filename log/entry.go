package log

import (
	"context"
	"net"
	"time"

	"github.com/phuslu/goid"
	"github.com/pubgo/funk/log/log_config"
	"github.com/pubgo/funk/log/log_fields"
	"github.com/pubgo/funk/logger"
	"github.com/pubgo/funk/stack"
)

func newEntry(log *loggerImpl, level logger.Level) logger.Entry {
	return &entryImpl{
		log:           log,
		level:         level,
		callerDepth:   log.callerDepth,
		callerEnabled: log.callerEnabled,
	}
}

var _ logger.Entry = (*entryImpl)(nil)

type entryImpl struct {
	log           *loggerImpl
	fields        logger.Fields
	level         logger.Level
	callerDepth   int
	callerEnabled bool
	ctx           context.Context
}

func (e *entryImpl) BytesLValue(key string, val logger.Valuer) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) RawLValue(key string, val logger.Valuer) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) BytesValue(key string, val logger.Valuer) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) RawValue(key string, val logger.Valuer) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) WithCaller() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) WithStack() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) WithXID() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) WithGoID() logger.Entry {
	e.fields = append(e.fields, log_fields.Int64("goid", goid.Goid()))
	return e
}

func (e *entryImpl) WithCtx(ctx context.Context) logger.Entry {
	e.ctx = ctx
	return e
}

func (e *entryImpl) GetFields() map[string]logger.Field {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) GetCtx() context.Context {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) Base64(dst, val []byte) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *entryImpl) GetName() string {
	return e.log.name
}

func (e *entryImpl) GetLevel() logger.Level {
	return e.level
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

func (e *entryImpl) Print(msg string) {
	e.fields = append(e.fields, log_fields.Time(log_config.FieldTime, time.Now().UTC()))
	e.fields = append(e.fields, log_fields.String(log_config.FieldMsgKey, msg))

	if e.callerEnabled {
		e.fields = append(e.fields, log_fields.String(log_config.FieldCallerKey, stack.Caller(e.callerDepth).String()))
	}

	var ent logger.Entry = e
	for _, h := range append(hooks, e.log.hooks...) {
		ent = h.Do(ent)
	}

	writer.WriteEntry(ent)
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
	e.fields = append(e.fields, log_fields.String(key, value))
	return e
}

func (e *entryImpl) StrL(key string, value ...string) logger.Entry {
	e.fields = append(e.fields, log_fields.StringL(key, value...))
	return e
}
