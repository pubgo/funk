package log

import (
	"context"
	"net"
	"time"

	"github.com/pubgo/funk/v2/log"
	"github.com/rs/zerolog"
)

var _ log.Event = (*eventImpl)(nil)

type eventImpl struct {
	evt *zerolog.Event
}

func (e eventImpl) Dict(key string, dict log.Event) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Enabled() bool {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Discard() log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Msg(msg string) {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Msgf(format string, v ...any) {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Fields(fields interface{}) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) MACAddr(key string, ha net.HardwareAddr) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) IPPrefix(key string, pfx net.IPNet) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) IPAddr(key string, ip net.IP) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Caller(skip ...int) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) CallerSkipFrame(skip int) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Type(key string, val interface{}) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Any(key string, i any) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Durs(key string, d []time.Duration) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Dur(key string, d time.Duration) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Times(key string, t []time.Time) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Time(key string, t time.Time) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Timestamp() log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Floats64(key string, f []float64) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Float64(key string, f float64) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Floats32(key string, f []float32) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Float32(key string, f float32) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Uints64(key string, i []uint64) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Uint64(key string, i uint64) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Err(err error) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Ctx(ctx context.Context) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Func(f func(e log.Event)) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Str(key string, value string) log.Event {
	//TODO implement me
	panic("implement me")
}

func (e eventImpl) Int(key string, value int) log.Event {
	//TODO implement me
	panic("implement me")
}
