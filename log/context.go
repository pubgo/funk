package log

import "context"

type ctxKey struct{}

func WithCtx(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// Ctx returns the Logger associated with the ctx.
func Ctx(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(*loggerImpl); ok && l != nil {
		return l
	}
	return stdLog
}
