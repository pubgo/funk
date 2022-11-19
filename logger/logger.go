package logger

import "context"

type Logger interface {
	Enabled() bool
	V(level int) Logger
	WithName(name string) Logger
	WithCtx(ctx context.Context) Logger
	WithValues(keysAndValues ...interface{}) Logger
	Info(level uint, msg string, kvs ...interface{}) error
	Error(err error, msg string, kvs ...interface{}) error
}

type CtxParser func(ctx context.Context)
type ValueParser func(v interface{}) (bool, string)
