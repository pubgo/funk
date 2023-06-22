package log

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

var _ Logger = (*loggerImpl)(nil)

type loggerImpl struct {
	name       string
	log        *zerolog.Logger
	fields     Map
	content    []byte
	callerSkip int
	lvl        Level
	filters    map[string]func(name string) bool
}

func (l *loggerImpl) WithNameFilter(filters map[string]func(name string) bool) Logger {
	var log = l.copy()
	log.filters = filters
	return log
}

func (l *loggerImpl) WithLevel(lvl Level) Logger {
	var log = l.copy()
	log.lvl = lvl
	return log
}

func (l *loggerImpl) WithEvent(evt *Event) Logger {
	if evt == nil {
		return l
	}

	evt1 := convertEvent(evt)
	if evt1.buf[0] == '{' && len(evt1.buf) == 1 {
		return l
	}

	var log = l.copy()
	logContent := make([]byte, 0, len(evt1.buf)+len(l.content))
	logContent = append(logContent, evt1.buf[1:]...)
	if len(log.content) > 0 {
		logContent = append(logContent, ',')
		logContent = append(logContent, log.content...)
	}
	log.content = logContent
	putEvent(evt)
	return log
}

func (l *loggerImpl) WithCallerSkip(skip int) Logger {
	if skip == 0 {
		return l
	}

	var log = l.copy()
	log.callerSkip += skip
	return log
}

func (l *loggerImpl) WithName(name string) Logger {
	if name == "" {
		return l
	}

	var log = l.copy()
	if log.name == "" {
		log.name = name
	} else {
		log.name = fmt.Sprintf("%s.%s", log.name, name)
	}
	return log
}

func (l *loggerImpl) WithFields(m Map) Logger {
	if len(m) == 0 {
		return l
	}

	var log = l.copy()
	var logFields = make(Map, len(m)+len(log.fields))
	for k, v := range m {
		logFields[k] = v
	}

	for k, v := range log.fields {
		logFields[k] = v
	}

	log.fields = logFields
	return log
}

func (l *loggerImpl) isFilterEnabled() bool {
	if l.filters == nil {
		return true
	}

	filter := l.filters[l.name]
	if filter == nil {
		return true
	}

	return filter(l.name)
}

func (l *loggerImpl) Debug(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.DebugLevel) {
		return nil
	}

	if !l.isFilterEnabled() {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Debug().Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) Info(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.InfoLevel) {
		return nil
	}

	if !l.isFilterEnabled() {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Info().Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) Warn(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.WarnLevel) {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Warn().Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) Error(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.ErrorLevel) {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Error().Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) Err(err error, ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.ErrorLevel) {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Err(err).Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) Panic(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.PanicLevel) {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Panic().Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) Fatal(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.FatalLevel) {
		return nil
	}

	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

	return l.newEvent(func(log *zerolog.Logger) *zerolog.Event { return log.Fatal().Func(withEventCtx(ctx)) })
}

func (l *loggerImpl) enabled(lvl zerolog.Level) bool {
	return lvl >= l.lvl && lvl >= zerolog.GlobalLevel()
}

func (l *loggerImpl) copy() *loggerImpl {
	var log = *l
	return &log
}

func (l *loggerImpl) getLog() *zerolog.Logger {
	if l.log != nil {
		return l.log
	}
	return stdZerolog
}

func (l *loggerImpl) newEvent(fn func(log *zerolog.Logger) *zerolog.Event) *zerolog.Event {
	var log = l.getLog()
	e := fn(log)
	if l.name != "" {
		e = e.Str("logger", l.name)
	}

	if l.callerSkip != 0 {
		e = e.CallerSkipFrame(l.callerSkip)
	}

	if l.fields != nil && len(l.fields) > 0 {
		e = e.Fields(l.fields)
	}

	if len(l.content) > 0 {
		e1 := convertEvent(e)
		e1.buf = append(e1.buf, ',')
		e1.buf = append(e1.buf, l.content...)
	}
	return e
}
