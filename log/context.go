package log

import "context"

type ctxKey struct{}

func WithCtx(ctx context.Context, log Logger) context.Context {
	if log == nil {
		panic("log is nil")
	}
	return context.WithValue(ctx, ctxKey{}, log)
}

// Ctx returns the Logger associated with the ctx.
func Ctx(ctx context.Context, defLog ...Logger) Logger {
	if l, ok := ctx.Value(ctxKey{}).(*loggerImpl); ok && l != nil {
		return l
	}

	if len(defLog) > 0 && defLog[0] != nil {
		return defLog[0]
	}

	return stdLog
}
