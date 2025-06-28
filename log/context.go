package log

import (
	"context"
)

type (
	ctxEventKey    struct{}
	ctxLoggerKey   struct{}
	disableLogKey  struct{}
	ctxMapFieldKey struct{}
)

func LoggerFromCtx(ctx context.Context, loggers ...Logger) Logger {
	defaultLog := stdLog
	if len(loggers) > 0 {
		defaultLog = loggers[0]
	}

	if ctx == nil {
		return defaultLog
	}

	if ll, ok := ctx.Value(ctxLoggerKey{}).(Logger); ok {
		return ll
	}

	return defaultLog
}

func CreateLoggerCtx(ctx context.Context, ll Logger) context.Context {
	if ll == nil || ctx == nil {
		panic("ctx or log param is nil")
	}

	return context.WithValue(ctx, ctxLoggerKey{}, ll)
}

func CreateEventCtx(ctx context.Context, evt *Event) context.Context {
	if evt == nil || ctx == nil {
		panic("ctx or log event is nil")
	}

	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func UpdateEventCtx(ctx context.Context, fields Map) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(fields) == 0 {
		return ctx
	}

	var evt = NewEvent()
	if e := getEventFromCtx(ctx); e != nil {
		evt = e
	} else {
		ctx = context.WithValue(ctx, ctxEventKey{}, evt)
	}

	for k, v := range fields {
		evt.Any(k, v)
	}

	return ctx
}

func getEventFromCtx(ctx context.Context) *Event {
	evt, ok := ctx.Value(ctxEventKey{}).(*Event)
	if ok {
		return evt
	}
	return nil
}

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

func createFieldCtx(ctx context.Context, mm Map) context.Context {
	if ctx == nil {
		panic("ctx is nil")
	}

	if len(mm) == 0 {
		return ctx
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
