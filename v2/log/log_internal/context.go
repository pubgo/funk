package log_internal

import (
	"context"
)

type ctxEventKey struct{}

func CreateEventCtx(ctx context.Context, evt Event) context.Context {
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

	e := GetEventFromCtx(ctx)
	if e != nil {
		for k, v := range fields {
			e.Any(k, v)
		}
	}

	return context.WithValue(ctx, ctxEventKey{}, e)
}

func GetEventFromCtx(ctx context.Context) Event {
	evt, ok := ctx.Value(ctxEventKey{}).(Event)
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

func IsLogDisabled(ctx context.Context) bool {
	b, ok := ctx.Value(disableLogKey{}).(bool)
	return b && ok
}
