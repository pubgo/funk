package log_internal

import (
	"context"
)

type ctxKey struct{}

func CreateCtx(ctx context.Context, evt Map) context.Context {
	if evt == nil || ctx == nil {
		panic("ctx or log event is nil")
	}

	return context.WithValue(ctx, ctxKey{}, evt)
}

func UpdateCtx(ctx context.Context, fields Map) context.Context {
	if ctx == nil {
		panic("ctx is nil")
	}

	if len(fields) == 0 {
		return ctx
	}

	for k, v := range GetFromCtx(ctx) {
		fields[k] = v
	}

	return context.WithValue(ctx, ctxKey{}, fields)
}

func GetFromCtx(ctx context.Context) Map {
	evt, ok := ctx.Value(ctxKey{}).(Map)
	if ok {
		return evt
	}
	return nil
}

type disableLogKey struct{}

func WithDisabled(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, disableLogKey{}, true)
}

func isLogDisabled(ctx context.Context) bool {
	b, ok := ctx.Value(disableLogKey{}).(bool)
	return b && ok
}
