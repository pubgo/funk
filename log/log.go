package log

import (
	"fmt"

	"github.com/rs/zerolog"
)

func New(log *zerolog.Logger) Logger {
	return &loggerImpl{
		log: log,
	}
}

var _ Logger = (*loggerImpl)(nil)

type loggerImpl struct {
	name       string
	log        *zerolog.Logger
	hooks      []zerolog.Hook
	context    Map
	callerSkip int
}

func (l *loggerImpl) WithCallerSkip(skip int) Logger {
	var log = l.copy()
	log.callerSkip += skip
	return log
}

func (l *loggerImpl) enabled(lvl zerolog.Level) bool {
	if lvl >= zerolog.GlobalLevel() {
		return true
	}
	return false
}

func (l *loggerImpl) copy() *loggerImpl {
	var log = *l
	return &log
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

func (l *loggerImpl) getLog() *zerolog.Logger {
	if l.log != nil {
		return l.log
	}

	return stdZero
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
		e = log.Error()
	case zerolog.WarnLevel:
		e = log.Warn()
	case zerolog.FatalLevel:
		e = log.Fatal()
	}

	if l.name != "" {
		e = e.Str("logger", l.name)
	}

	if l.callerSkip != 0 {
		e = e.CallerSkipFrame(l.callerSkip)
	}

	if l.context != nil && len(l.context) > 0 {
		e = e.Fields(l.context)
	}

	return e
}

func (l *loggerImpl) WithFields(m Map) Logger {
	var log = l.copy()
	log.context = m
	return log
}

func (l *loggerImpl) WithHooks(hooks ...zerolog.Hook) Logger {
	var log = l.copy()
	log.hooks = append(log.hooks, hooks...)
	return log
}

func (l *loggerImpl) WithCaller(depth int) Logger {
	var log = l.copy()
	log.callerSkip += depth
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
