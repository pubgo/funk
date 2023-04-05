package errors

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
)

type EventEnc struct {
	evt *zerolog.Event
}

func (e *EventEnc) Msg(msg string) {
	e.evt.Msg(msg)
}

func (e *EventEnc) Msgf(format string, v ...interface{}) {
	e.evt.Msgf(format, v...)
}

func (e *EventEnc) MsgFunc(createMsg func() string) {
	e.evt.MsgFunc(createMsg)
}

func (e *EventEnc) Fields(fields interface{}) *EventEnc {
	e.evt.Fields(fields)
	return e
}

// Func allows an anonymous func to run only if the event is enabled.
func (e *EventEnc) Func(f func(e *Event)) *EventEnc {
	e.evt.Func(f)
	return e
}

func (e *EventEnc) Str(key, val string) *EventEnc {
	e.evt.Str(key, val)
	return e
}

func (e *EventEnc) Strs(key string, vals []string) *EventEnc {
	e.evt.Strs(key, vals)
	return e
}

func (e *EventEnc) Stringer(key string, val fmt.Stringer) *EventEnc {
	e.evt.Stringer(key, val)
	return e
}

func (e *EventEnc) Stringers(key string, vals []fmt.Stringer) *EventEnc {
	e.evt.Stringers(key, vals)
	return e
}

func (e *EventEnc) Bytes(key string, val []byte) *EventEnc {
	e.evt.Bytes(key, val)
	return e
}

func (e *EventEnc) Hex(key string, val []byte) *EventEnc {
	e.evt.Hex(key, val)
	return e
}

func (e *EventEnc) RawJSON(key string, b []byte) *EventEnc {
	e.evt.RawJSON(key, b)
	return e
}

func (e *EventEnc) AnErr(key string, err error) *EventEnc {
	e.evt.AnErr(key, err)
	return e
}

func (e *EventEnc) Errs(key string, errs []error) *EventEnc {
	e.evt.Errs(key, errs)
	return e
}

func (e *EventEnc) Err(err error) *EventEnc {
	e.evt.Err(err)
	return e
}

func (e *EventEnc) Bool(key string, b bool) *EventEnc {
	e.evt.Bool(key, b)
	return e
}

func (e *EventEnc) Bools(key string, b []bool) *EventEnc {
	e.evt.Bools(key, b)
	return e
}

func (e *EventEnc) Int(key string, i int) *EventEnc {
	e.evt.Int(key, i)
	return e
}

func (e *EventEnc) Ints(key string, i []int) *EventEnc {
	e.evt.Ints(key, i)
	return e
}

func (e *EventEnc) Int8(key string, i int8) *EventEnc {
	e.evt.Int8(key, i)
	return e
}

func (e *EventEnc) Ints8(key string, i []int8) *EventEnc {
	e.evt.Ints8(key, i)
	return e
}

func (e *EventEnc) Int16(key string, i int16) *EventEnc {
	e.evt.Int16(key, i)
	return e
}

func (e *EventEnc) Ints16(key string, i []int16) *EventEnc {
	e.evt.Ints16(key, i)
	return e
}

func (e *EventEnc) Int32(key string, i int32) *EventEnc {
	e.evt.Int32(key, i)
	return e
}

func (e *EventEnc) Ints32(key string, i []int32) *EventEnc {
	e.evt.Ints32(key, i)
	return e
}

func (e *EventEnc) Int64(key string, i int64) *EventEnc {
	e.evt.Int64(key, i)
	return e
}

func (e *EventEnc) Ints64(key string, i []int64) *EventEnc {
	e.evt.Ints64(key, i)
	return e
}

func (e *EventEnc) Uint(key string, i uint) *EventEnc {
	e.evt.Uint(key, i)
	return e
}

func (e *EventEnc) Uints(key string, i []uint) *EventEnc {
	e.evt.Uints(key, i)
	return e
}

func (e *EventEnc) Uint8(key string, i uint8) *EventEnc {
	e.evt.Uint8(key, i)
	return e
}

func (e *EventEnc) Uints8(key string, i []uint8) *EventEnc {
	e.evt.Uints8(key, i)
	return e
}

func (e *EventEnc) Uint16(key string, i uint16) *EventEnc {
	e.evt.Uint16(key, i)
	return e
}

func (e *EventEnc) Uints16(key string, i []uint16) *EventEnc {
	e.evt.Uints16(key, i)
	return e
}

func (e *EventEnc) Uint32(key string, i uint32) *EventEnc {
	e.evt.Uint32(key, i)
	return e
}

func (e *EventEnc) Uints32(key string, i []uint32) *EventEnc {
	e.evt.Uints32(key, i)
	return e
}

func (e *EventEnc) Uint64(key string, i uint64) *EventEnc {
	e.evt.Uint64(key, i)
	return e
}

func (e *EventEnc) Uints64(key string, i []uint64) *EventEnc {
	e.evt.Uints64(key, i)
	return e
}

func (e *EventEnc) Float32(key string, f float32) *EventEnc {
	e.evt.Float32(key, f)
	return e
}

func (e *EventEnc) Floats32(key string, f []float32) *EventEnc {
	e.evt.Floats32(key, f)
	return e
}

func (e *EventEnc) Float64(key string, f float64) *EventEnc {
	e.evt.Float64(key, f)
	return e
}

func (e *EventEnc) Floats64(key string, f []float64) *EventEnc {
	e.evt.Floats64(key, f)
	return e
}

func (e *EventEnc) Timestamp() *EventEnc {
	e.evt.Timestamp()
	return e
}

func (e *EventEnc) Time(key string, t time.Time) *EventEnc {
	e.evt.Time(key, t)
	return e
}

func (e *EventEnc) Times(key string, t []time.Time) *EventEnc {
	e.evt.Times(key, t)
	return e
}

func (e *EventEnc) Dur(key string, d time.Duration) *EventEnc {
	e.evt.Dur(key, d)
	return e
}

func (e *EventEnc) Durs(key string, d []time.Duration) *EventEnc {
	e.evt.Durs(key, d)
	return e
}

func (e *EventEnc) TimeDiff(key string, t time.Time, start time.Time) *EventEnc {
	e.evt.TimeDiff(key, t, start)
	return e
}

// Any is a wrapper around Event.Interface.
func (e *EventEnc) Any(key string, i interface{}) *EventEnc {
	e.evt.Any(key, i)
	return e
}

func (e *EventEnc) Interface(key string, i interface{}) *EventEnc {
	e.evt.Interface(key, i)
	return e
}

func (e *EventEnc) Type(key string, val interface{}) *EventEnc {
	e.evt.Type(key, val)
	return e
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *EventEnc) IPAddr(key string, ip net.IP) *EventEnc {
	e.evt.IPAddr(key, ip)
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *EventEnc) IPPrefix(key string, pfx net.IPNet) *EventEnc {
	e.evt.IPPrefix(key, pfx)
	return e
}

// MACAddr adds MAC address to the event
func (e *EventEnc) MACAddr(key string, ha net.HardwareAddr) *EventEnc {
	e.evt.MACAddr(key, ha)
	return e
}
