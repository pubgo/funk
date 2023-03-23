package log

import "context"

type ctxKey struct{}

func (l *loggerImpl) ParseCtx(ctx context.Context) Logger {

}

func (l *loggerImpl) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// Ctx returns the Logger associated with the ctx.
func Ctx(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(*loggerImpl); ok && l != nil {
		return l
	}
	return stdLog
}
