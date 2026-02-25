package log

import (
	"context"
	"log"
	"maps"
)

type (
	ctxEventKey    struct{}
	ctxLoggerKey   struct{}
	disableLogKey  struct{}
	ctxMapFieldKey struct{}
)

func GetFromCtx(ctx context.Context, loggers ...Logger) Logger {
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

func CreateCtx(ctx context.Context, ll Logger) context.Context {
	if ll == nil || ctx == nil {
		log.Panicln("ctx or log param is nil")
	}

	return context.WithValue(ctx, ctxLoggerKey{}, ll)
}

func CreateFieldsCtx(ctx context.Context, evt Fields) context.Context {
	if evt == nil || ctx == nil {
		log.Panicln("ctx or log event is nil")
	}

	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func UpdateFieldsCtx(ctx context.Context, fields Fields) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(fields) == 0 {
		return ctx
	}

	evt := make(Fields)
	if e := GetFieldsFromCtx(ctx); e != nil {
		evt = e
	}

	maps.Copy(evt, fields)
	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func GetFieldsFromCtx(ctx context.Context) Fields {
	evt, ok := ctx.Value(ctxEventKey{}).(Fields)
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

type fieldMap struct {
	fields Fields
	name   string
}

func createFieldCtx(ctx context.Context, field *fieldMap) context.Context {
	if ctx == nil {
		panic("ctx is nil")
	}

	return context.WithValue(ctx, ctxMapFieldKey{}, field)
}

func getFieldFromCtx(ctx context.Context) *fieldMap {
	if ctx == nil {
		return nil
	}

	field, ok := ctx.Value(ctxMapFieldKey{}).(*fieldMap)
	if ok {
		return field
	}
	return nil
}
