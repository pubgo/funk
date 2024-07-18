package log

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pubgo/funk/stack"
	"github.com/rs/zerolog"
)

var _ Logger = (*loggerImpl)(nil)

type loggerImpl struct {
	name       string
	log        *zerolog.Logger
	fields     Map
	content    *Event
	callerSkip int
	lvl        Level
}

func (l *loggerImpl) WithLevel(lvl Level) Logger {
	log := l.copy()
	log.lvl = lvl
	return log
}

func (l *loggerImpl) WithEvent(evt *Event) Logger {
	if evt == nil {
		return l
	}

	log := l.copy()
	log.content = mergeEvent(l.content, evt)

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

func (l *loggerImpl) WithName(name string) Logger {
	name = strings.TrimSpace(name)
	if name == "" {
		return l
	}

	log := l.copy()
	if log.fields == nil {
		log.fields = Map{ModuleName: stack.Caller(1).Pkg}
	} else {
		log.fields[ModuleName] = stack.Caller(1).Pkg
	}

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

	log := l.copy()
	logFields := make(Map, len(m)+len(log.fields))
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
	if len(ctxL) > 0 {
		ctx = ctxL[0]
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

	if err != nil {
		if errJson, ok := err.(json.Marshaler); ok {
			errJsonBytes, _ := errJson.MarshalJSON()
			if len(errJsonBytes) > 0 {
				return l.newEvent(ctx, l.getLog().Error().Str("error", err.Error()).RawJSON("error_detail", errJsonBytes))
			}
		}

		return l.newEvent(ctx, l.getLog().Error().Str("error", err.Error()))
	}

	return l.newEvent(ctx, l.getLog().Err(err))
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

	enabled := true
	if logEnableChecker != nil {
		enabled = logEnableChecker(ctx, lvl, l.name, l.fields)
	}
	return enabled && lvl >= l.lvl && lvl >= zerolog.GlobalLevel()
}

func (l *loggerImpl) copy() *loggerImpl {
	log := *l
	return &log
}

func (l *loggerImpl) getLog() *zerolog.Logger {
	if l.log != nil {
		return l.log
	}
	return stdZeroLog
}

func (l *loggerImpl) newEvent(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	if l.name != "" {
		e = e.Str("logger", l.name)
	}

	if l.callerSkip != 0 {
		e = e.CallerSkipFrame(l.callerSkip)
	}

	if l.fields != nil && len(l.fields) > 0 {
		e = e.Fields(l.fields)
	}

	return mergeEvent(e, getEventFromCtx(ctx), l.content)
}
