package log_internal

import (
	"context"
	"encoding/json"
	"net"
	"time"
)

var _ Event = &noopEvent{}

type noopEvent struct{}

func (n *noopEvent) Enabled() bool                                   { return false }
func (n *noopEvent) Msg(msg string)                                  {}
func (n *noopEvent) Msgf(format string, v ...any)                    {}
func (n *noopEvent) Dict(key string, dict Event) Event               { return n }
func (n *noopEvent) Discard() Event                                  { return n }
func (n *noopEvent) Fields(fields interface{}) Event                 { return n }
func (n *noopEvent) MACAddr(key string, ha net.HardwareAddr) Event   { return n }
func (n *noopEvent) IPPrefix(key string, pfx net.IPNet) Event        { return n }
func (n *noopEvent) IPAddr(key string, ip net.IP) Event              { return n }
func (n *noopEvent) Caller(skip ...int) Event                        { return n }
func (n *noopEvent) CallerSkipFrame(skip int) Event                  { return n }
func (n *noopEvent) Type(key string, val interface{}) Event          { return n }
func (n *noopEvent) Any(key string, i any) Event                     { return n }
func (n *noopEvent) Durs(key string, d []time.Duration) Event        { return n }
func (n *noopEvent) Dur(key string, d time.Duration) Event           { return n }
func (n *noopEvent) Times(key string, t []time.Time) Event           { return n }
func (n *noopEvent) Time(key string, t time.Time) Event              { return n }
func (n *noopEvent) Timestamp() Event                                { return n }
func (n *noopEvent) Floats64(key string, f []float64) Event          { return n }
func (n *noopEvent) Float64(key string, f float64) Event             { return n }
func (n *noopEvent) Floats32(key string, f []float32) Event          { return n }
func (n *noopEvent) Float32(key string, f float32) Event             { return n }
func (n *noopEvent) Uints64(key string, i []uint64) Event            { return n }
func (n *noopEvent) Uint64(key string, i uint64) Event               { return n }
func (n *noopEvent) Err(err error) Event                             { return n }
func (n *noopEvent) Ctx(ctx context.Context) Event                   { return n }
func (n *noopEvent) Func(f func(e Event)) Event                      { return n }
func (n *noopEvent) Str(key string, value string) Event              { return n }
func (n *noopEvent) Int(key string, value int) Event                 { return n }
func (n *noopEvent) Bool(key string, value bool) Event               { return n }
func (n *noopEvent) RawJSON(key string, value json.RawMessage) Event { return n }

var _ EventLogger = &noopEventLogger{}

type noopEventLogger struct{}

func (n noopEventLogger) Debug(ctx ...context.Context) Event          { return new(noopEvent) }
func (n noopEventLogger) Info(ctx ...context.Context) Event           { return new(noopEvent) }
func (n noopEventLogger) Warn(ctx ...context.Context) Event           { return new(noopEvent) }
func (n noopEventLogger) Error(ctx ...context.Context) Event          { return new(noopEvent) }
func (n noopEventLogger) Err(err error, ctx ...context.Context) Event { return new(noopEvent) }
func (n noopEventLogger) Panic(ctx ...context.Context) Event          { return new(noopEvent) }
func (n noopEventLogger) Fatal(ctx ...context.Context) Event          { return new(noopEvent) }
