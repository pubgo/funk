package log

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

var _ Logger = (*loggerImpl)(nil)

type loggerImpl struct {
	name          string
	log           *zerolog.Logger
	fields        Map
	content       *Event
	callerSkip    int
	lvl           Level
	enableChecker func(lvl Level, name string, fields Map) bool
}

func (l *loggerImpl) WithEnableChecker(cb func(lvl Level, name string, fields Map) bool) Logger {
	var log = l.copy()
	log.enableChecker = cb
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

	var log = l.copy()
	log.content = mergeEvent(l.content, evt)

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
	if m == nil || len(m) == 0 {
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

func (l *loggerImpl) Debug(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.DebugLevel) {
		return nil
	}

	return l.newEvent(ctxL, l.getLog().Debug())
}

func (l *loggerImpl) Info(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.InfoLevel) {
		return nil
	}

	return l.newEvent(ctxL, l.getLog().Info())
}

func (l *loggerImpl) Warn(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.WarnLevel) {
		return nil
	}

	return l.newEvent(ctxL, l.getLog().Warn())
}

func (l *loggerImpl) Error(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.ErrorLevel) {
		return nil
	}

	return l.newEvent(ctxL, l.getLog().Error())
}

func (l *loggerImpl) Err(err error, ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.ErrorLevel) {
		return nil
	}

	if err != nil {
		if errJson, ok := err.(json.Marshaler); ok {
			var errJsonBytes, _ = errJson.MarshalJSON()
			if errJsonBytes != nil && len(errJsonBytes) > 0 {
				return l.newEvent(ctxL, l.getLog().Error().Str("error", err.Error()).RawJSON("error_detail", errJsonBytes))
			}
		}

		return l.newEvent(ctxL, l.getLog().Error().Str("error", err.Error()))
	}

	return l.newEvent(ctxL, l.getLog().Err(err))
}

func (l *loggerImpl) Panic(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.PanicLevel) {
		return nil
	}

	return l.newEvent(ctxL, l.getLog().Panic())
}

func (l *loggerImpl) Fatal(ctxL ...context.Context) *zerolog.Event {
	if !l.enabled(zerolog.FatalLevel) {
		return nil
	}

	return l.newEvent(ctxL, l.getLog().Fatal())
}

func (l *loggerImpl) enabled(lvl zerolog.Level) bool {
	if l.enableChecker != nil {
		if l.enableChecker(lvl, l.name, l.fields) {
			return lvl >= l.lvl && lvl >= zerolog.GlobalLevel()
		}
	}

	return false
}

func (l *loggerImpl) copy() *loggerImpl {
	var log = *l
	return &log
}

func (l *loggerImpl) getLog() *zerolog.Logger {
	if l.log != nil {
		return l.log
	}
	return stdZeroLog
}

func (l *loggerImpl) newEvent(ctxL []context.Context, e *zerolog.Event) *zerolog.Event {
	var ctx = context.Background()
	if len(ctxL) > 0 {
		ctx = ctxL[0]
	}

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
