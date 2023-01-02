package log

import "context"

type ctxKey struct{}

func (l *loggerImpl) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, &l)
}

// Ctx returns the Logger associated with the ctx.
func Ctx(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return stdLog
}
