package log

import "context"

type ctxEventKey struct{}

func CreateEventCtx(ctx context.Context, evt *Event) context.Context {
	if evt == nil || ctx == nil {
		panic("ctx or log event is nil")
	}

	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func withEventCtx(ctx context.Context) func(e *Event) {
	return func(e *Event) {
		var evt, ok = ctx.Value(ctxEventKey{}).(*Event)
		if ok {
			WithEvent(evt)(e)
		}
	}
}
