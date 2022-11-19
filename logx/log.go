package logx

import (
	"fmt"

	"github.com/go-logr/logr"
)

func WithCallDepth(depth int) logr.Logger {
	return gl.WithCallDepth(depth)
}

func WithName(name string) logr.Logger {
	return gl.WithName(name)
}

func V(level int) logr.Logger {
	return gl.V(level)
}

func WithValues(keysAndValues ...interface{}) logr.Logger {
	return gl.WithValues(keysAndValues...)
}

func IfEnabled(level int, fn func(log logr.Logger)) {
	var log = V(level)
	if log.Enabled() {
		fn(log)
	}
}

func Enabled() bool {
	return gl.Enabled()
}

func Info(msg string, keysAndValues ...interface{}) {
	gl.WithCallDepth(1).Info(msg, keysAndValues...)
}

func Infof(format string, args ...interface{}) {
	gl.WithCallDepth(1).Info(fmt.Sprintf(format, args...))
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	if err == nil {
		return
	}

	gl.WithCallDepth(1).Error(err, msg, keysAndValues...)
}
