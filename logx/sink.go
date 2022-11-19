package logx

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/DataDog/gostackparse"
	_ "github.com/DataDog/gostackparse"
	logkit "github.com/go-kit/log"
	"github.com/go-logr/logr"
	"github.com/pubgo/funk/assert"
)

var _ logr.LogSink = (*sink)(nil)
var _ logr.CallDepthLogSink = (*sink)(nil)

type sink struct {
	level     int
	callDepth int32
	prefix    string
	values    []interface{}
}

// Enabled reports whether this Logger is enabled with respect to the current global log level.
func (s sink) Enabled(level int) bool {
	if level > int(atomic.LoadInt32(&gv)) {
		return false
	}

	if s.prefix == "" {
		return true
	}

	return true
}

func (s sink) WithCallDepth(depth int) logr.LogSink {
	s.callDepth += int32(depth)
	return s
}

func (s sink) Init(info logr.RuntimeInfo) {
	atomic.AddInt32(&s.callDepth, int32(info.CallDepth))
}

func (s sink) Info(level int, msg string, keysAndValues ...interface{}) {
	if !s.Enabled(level) {
		return
	}

	keysAndValues = append(keysAndValues, s.values...)
	keysAndValues = append(keysAndValues, "caller", logkit.Caller(int(s.callDepth)+DefaultCallerSkip)())
	keysAndValues = append(keysAndValues, "logger", s.prefix)
	keysAndValues = append(keysAndValues, "level", "info")
	keysAndValues = append(keysAndValues, "msg", msg)
	keysAndValues = append(keysAndValues, "ts", time.Now().UTC().Format(TimestampFormat))
	assert.Must(logWriter.Info(keysAndValues...))
}

func (s sink) Error(err error, msg string, keysAndValues ...interface{}) {
	if err == nil || reflect.ValueOf(err).IsZero() {
		return
	}

	goroutines, _ := gostackparse.Parse(bytes.NewReader(debug.Stack()))
	keysAndValues = append(keysAndValues, s.values...)
	keysAndValues = append(keysAndValues, "caller", logkit.Caller(int(s.callDepth)+DefaultCallerSkip)())
	keysAndValues = append(keysAndValues, "logger", s.prefix)
	keysAndValues = append(keysAndValues, "level", "error")
	keysAndValues = append(keysAndValues, "msg", msg)
	keysAndValues = append(keysAndValues, "error", err.Error())
	keysAndValues = append(keysAndValues, "error_detail", fmt.Sprintf("%#v", err))
	keysAndValues = append(keysAndValues, "stacktrace", goroutines)
	keysAndValues = append(keysAndValues, "ts", time.Now().UTC().Format(TimestampFormat))
	assert.Must(logWriter.Error(err, msg, keysAndValues...))
}

func (s sink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	s.values = append(s.values, keysAndValues...)
	return s
}

func (s sink) WithName(name string) logr.LogSink {
	if len(s.prefix) > 0 {
		s.prefix = s.prefix + "."
	}
	s.prefix += name
	return s
}
