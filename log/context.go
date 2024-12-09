package log

import (
	"context"
)

type ctxEventKey struct{}

func CreateEventCtx(ctx context.Context, evt *Event) context.Context {
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

	var evt = NewEvent()
	for k, v := range fields {
		evt.Any(k, v)
	}

	if e := getEventFromCtx(ctx); e != nil {

		evt = mergeEvent(evt, e)
	}

	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func getEventFromCtx(ctx context.Context) *Event {
	evt, ok := ctx.Value(ctxEventKey{}).(*Event)
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

type ctxMapFieldKey struct{}

func createFieldCtx(ctx context.Context, mm Map) context.Context {
	if mm == nil || ctx == nil {
		panic("ctx or log field is nil")
	}

	return context.WithValue(ctx, ctxMapFieldKey{}, mm)
}

func getFieldFromCtx(ctx context.Context) Map {
	if ctx == nil {
		return make(Map)
	}
	
	field, ok := ctx.Value(ctxMapFieldKey{}).(Map)
	if ok {
		return field
	}
	return make(Map)
}
