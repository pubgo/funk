package log

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log/logfields"
)

var _ Logger = (*loggerImpl)(nil)

func New(log *zerolog.Logger) Logger {
	return &loggerImpl{
		log: log,
	}
}

type loggerImpl struct {
	name       string
	log        *zerolog.Logger
	fields     Fields
	callerSkip int
	lvl        Level
}

func (l *loggerImpl) WithLevel(lvl Level) Logger {
	log := l.copy()
	log.lvl = lvl
	return log
}

func (l *loggerImpl) WithCallerSkip(skip int) Logger {
	if skip == 0 {
		return l
	}

	log := l.copy()
	log.callerSkip += skip
	return log
}

func (l *loggerImpl) nameWithCaller(name string, caller int) Logger {
	name = strings.TrimSpace(name)
	if name == "" {
		return l
	}

	log := l.copy()
	if log.fields == nil {
		log.fields = make(Fields, 1)
	}

	if log.name == "" {
		log.name = name
	} else {
		log.name = fmt.Sprintf("%s.%s", log.name, name)
	}

	log.callerSkip += caller
	return log
}

func (l *loggerImpl) WithName(name string) Logger {
	return l.nameWithCaller(name, 0)
}

func (l *loggerImpl) WithFields(m Fields) Logger {
	if len(m) == 0 {
		return l
	}

	log := l.copy()
	logFields := make(Fields, len(m)+len(log.fields))
	for k, v := range m {
		logFields[k] = v
	}

	for k, v := range log.fields {
		logFields[k] = v
	}

	log.fields = logFields
	return log
}

func (l *loggerImpl) getCtx(ctxL ...context.Context) context.Context {
	ctx := context.Background()
	for i := range ctxL {
		if ctxL[i] != nil {
			ctx = ctxL[i]
			break
		}
	}
	return ctx
}

func (l *loggerImpl) Debug(ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.DebugLevel) {
		return nil
	}

	return l.newEvent(ctx, l.getLog().Debug())
}

func (l *loggerImpl) Info(ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.InfoLevel) {
		return nil
	}

	return l.newEvent(ctx, l.getLog().Info())
}

func (l *loggerImpl) Warn(ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.WarnLevel) {
		return nil
	}

	return l.newEvent(ctx, l.getLog().Warn())
}

func (l *loggerImpl) Error(ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.ErrorLevel) {
		return nil
	}

	return l.newEvent(ctx, l.getLog().Error())
}

func (l *loggerImpl) Err(err error, ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.ErrorLevel) {
		return nil
	}

	fn := func(e *zerolog.Event) {
		if err == nil {
			return
		}

		if id := errors.GetErrorId(err); id != "" {
			e.Str("error_id", id)
		}

		e.Str("error_detail", errDetail(err))
		e.Str(zerolog.ErrorFieldName, err.Error())
	}
	return l.newEvent(ctx, l.getLog().Err(err).Func(fn))
}

func (l *loggerImpl) Panic(ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.PanicLevel) {
		return nil
	}

	return l.newEvent(ctx, l.getLog().Panic())
}

func (l *loggerImpl) Fatal(ctxL ...context.Context) *zerolog.Event {
	ctx := l.getCtx(ctxL...)
	if !l.enabled(ctx, zerolog.FatalLevel) {
		return nil
	}

	return l.newEvent(ctx, l.getLog().Fatal())
}

func (l *loggerImpl) enabled(ctx context.Context, lvl zerolog.Level) bool {
	if isLogDisabled(ctx) {
		return false
	}

	return lvl >= l.lvl && lvl >= zerolog.GlobalLevel()
}

func (l *loggerImpl) copy() *loggerImpl {
	return &loggerImpl{
		log:        l.log,
		fields:     maps.Clone(l.fields),
		lvl:        l.lvl,
		name:       l.name,
		callerSkip: l.callerSkip,
	}
}

func (l *loggerImpl) getLog() *zerolog.Logger {
	if l.log != nil {
		return l.log
	}
	return stdZeroLog
}

func (l *loggerImpl) newEvent(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	name := l.name
	fields := l.fields

	if m, ok := fields[logfields.Module].(string); ok {
		name = m
	}

	if name != "" {
		e = e.Str(logfields.Logger, name)
	}

	if l.callerSkip != 0 {
		e = e.CallerSkipFrame(l.callerSkip)
	}

	if fields == nil {
		fields = GetFieldsFromCtx(ctx)
	} else {
		for k, v := range GetFieldsFromCtx(ctx) {
			fields[k] = v
		}
	}

	if len(fields) > 0 {
		e = e.Fields(fields)
	}

	e = e.Ctx(createFieldCtx(ctx, &fieldMap{name: name, fields: fields}))
	return e
}
