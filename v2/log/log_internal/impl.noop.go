package log_internal

import (
	"context"
	"encoding/json"
	"net"
	"time"
)

var _ Event = &noopLogger{}

type noopLogger struct{}

func (n *noopLogger) Enabled() bool                                   { return false }
func (n *noopLogger) Msg(msg string)                                  {}
func (n *noopLogger) Msgf(format string, v ...any)                    {}
func (n *noopLogger) Dict(key string, dict Event) Event               { return n }
func (n *noopLogger) Discard() Event                                  { return n }
func (n *noopLogger) Fields(fields interface{}) Event                 { return n }
func (n *noopLogger) MACAddr(key string, ha net.HardwareAddr) Event   { return n }
func (n *noopLogger) IPPrefix(key string, pfx net.IPNet) Event        { return n }
func (n *noopLogger) IPAddr(key string, ip net.IP) Event              { return n }
func (n *noopLogger) Caller(skip ...int) Event                        { return n }
func (n *noopLogger) CallerSkipFrame(skip int) Event                  { return n }
func (n *noopLogger) Type(key string, val interface{}) Event          { return n }
func (n *noopLogger) Any(key string, i any) Event                     { return n }
func (n *noopLogger) Durs(key string, d []time.Duration) Event        { return n }
func (n *noopLogger) Dur(key string, d time.Duration) Event           { return n }
func (n *noopLogger) Times(key string, t []time.Time) Event           { return n }
func (n *noopLogger) Time(key string, t time.Time) Event              { return n }
func (n *noopLogger) Timestamp() Event                                { return n }
func (n *noopLogger) Floats64(key string, f []float64) Event          { return n }
func (n *noopLogger) Float64(key string, f float64) Event             { return n }
func (n *noopLogger) Floats32(key string, f []float32) Event          { return n }
func (n *noopLogger) Float32(key string, f float32) Event             { return n }
func (n *noopLogger) Uints64(key string, i []uint64) Event            { return n }
func (n *noopLogger) Uint64(key string, i uint64) Event               { return n }
func (n *noopLogger) Err(err error) Event                             { return n }
func (n *noopLogger) Ctx(ctx context.Context) Event                   { return n }
func (n *noopLogger) Func(f func(e Event)) Event                      { return n }
func (n *noopLogger) Str(key string, value string) Event              { return n }
func (n *noopLogger) Int(key string, value int) Event                 { return n }
func (n *noopLogger) Bool(key string, value bool) Event               { return n }
func (n *noopLogger) RawJSON(key string, value json.RawMessage) Event { return n }
