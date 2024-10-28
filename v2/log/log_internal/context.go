package log_internal

import (
	"context"
)

type ctxEventKey struct{}

func CreateEventCtx(ctx context.Context, evt Map) context.Context {
	if evt == nil || ctx == nil {
		panic("ctx or log event is nil")
	}

	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func UpdateEventCtx(ctx context.Context, fields Map) context.Context {
	if ctx == nil {
		panic("ctx is nil")
	}

	if len(fields) == 0 {
		return ctx
	}

	for k, v := range GetEventFromCtx(ctx) {
		fields[k] = v
	}

	return context.WithValue(ctx, ctxEventKey{}, fields)
}

func GetEventFromCtx(ctx context.Context) Map {
	evt, ok := ctx.Value(ctxEventKey{}).(Map)
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
