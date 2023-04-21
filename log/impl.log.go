package log

import (
	"fmt"

	"github.com/rs/zerolog"
)

var _ Logger = (*loggerImpl)(nil)

type loggerImpl struct {
	name       string
	log        *zerolog.Logger
	hooks      []zerolog.Hook
	fields     Map
	content    []byte
	callerSkip int
	lvl        Level
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

func (l *loggerImpl) WithHooks(hooks ...zerolog.Hook) Logger {
	if len(hooks) == 0 {
		return l
	}

	var log = l.copy()
	var logHook = make([]zerolog.Hook, 0, len(hooks)+len(log.hooks))
	logHook = append(logHook, hooks...)
	logHook = append(logHook, log.hooks...)
	log.hooks = logHook
	return log
}

func (l *loggerImpl) Debug() *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.DebugLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.DebugLevel)
}

func (l *loggerImpl) Info() *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.InfoLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.InfoLevel)
}

func (l *loggerImpl) Warn() *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.WarnLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.WarnLevel)
}

func (l *loggerImpl) Error() *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.ErrorLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.ErrorLevel)
}

func (l *loggerImpl) Err(err error) *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.ErrorLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.ErrorLevel).Err(err)
}

func (l *loggerImpl) Panic() *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.PanicLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.PanicLevel)
}

func (l *loggerImpl) Fatal() *zerolog.Event {
	var log = l.getLog()
	if !l.enabled(zerolog.FatalLevel) {
		return nil
	}

	return l.newEvent(*log, zerolog.FatalLevel)
}

func (l *loggerImpl) Print(msg string) {
	var log = l.getLog()
	if !l.enabled(zerolog.DebugLevel) {
		return
	}

	l.newEvent(*log, zerolog.DebugLevel).CallerSkipFrame(1).Msg(msg)
}

func (l *loggerImpl) Printf(format string, args ...any) {
	var log = l.getLog()
	if !l.enabled(zerolog.DebugLevel) {
		return
	}

	l.newEvent(*log, zerolog.DebugLevel).CallerSkipFrame(1).Msgf(format, args...)
}

func (l *loggerImpl) enabled(lvl zerolog.Level) bool {
	return lvl >= l.lvl
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

func (l *loggerImpl) newEvent(log zerolog.Logger, level zerolog.Level) *zerolog.Event {
	for i := range l.hooks {
		log = log.Hook(l.hooks[i])
	}

	var e *zerolog.Event
	switch level {
	case zerolog.DebugLevel:
		e = log.Debug()
	case zerolog.InfoLevel:
		e = log.Info()
	case zerolog.ErrorLevel:
		e = log.Error()
	case zerolog.PanicLevel:
		e = log.Panic()
	case zerolog.WarnLevel:
		e = log.Warn()
	case zerolog.FatalLevel:
		e = log.Fatal()
	default:
		e = log.Debug()
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

	if len(l.content) > 0 {
		e1 := convertEvent(e)
		e1.buf = append(e1.buf, ',')
		e1.buf = append(e1.buf, l.content...)
	}
	return e
}
