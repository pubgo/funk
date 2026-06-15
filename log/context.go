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

// FromCtx returns the logger stored in the context or the global logger when
// the context does not carry one.
func FromCtx(ctx context.Context) Logger {
	return GetFromCtx(ctx)
}

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

// WithLogger attaches a logger to the context and is safe to use with nil input.
// When ctx is nil it falls back to context.Background(); when ll is nil it falls
// back to the global logger.
func WithLogger(ctx context.Context, ll Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if ll == nil {
		ll = stdLog
	}

	return context.WithValue(ctx, ctxLoggerKey{}, ll)
}

func CreateFieldsCtx(ctx context.Context, evt Fields) context.Context {
	if evt == nil || ctx == nil {
		log.Panicln("ctx or log event is nil")
	}

	return context.WithValue(ctx, ctxEventKey{}, maps.Clone(evt))
}

// WithFields adds or overrides fields in the context used by log events.
// It is a nil-safe convenience wrapper around UpdateFieldsCtx.
func WithFields(ctx context.Context, fields Fields) context.Context {
	return UpdateFieldsCtx(ctx, fields)
}

func UpdateFieldsCtx(ctx context.Context, fields Fields) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(fields) == 0 {
		return ctx
	}

	evt := maps.Clone(GetFieldsFromCtx(ctx))
	if evt == nil {
		evt = make(Fields, len(fields))
	}

	maps.Copy(evt, fields)
	return context.WithValue(ctx, ctxEventKey{}, evt)
}

func GetFieldsFromCtx(ctx context.Context) Fields {
	if ctx == nil {
		return nil
	}

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
	if ctx == nil {
		return false
	}

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
