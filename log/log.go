package log

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/logger"
)

func New() logger.Logger {
	return &loggerImpl{}
}

var _ logger.Logger = (*loggerImpl)(nil)

type loggerImpl struct {
	name          string
	callerDepth   int
	callerEnabled bool
	fields        []logger.Field
	hooks         []logger.Hook
}

func (l *loggerImpl) GetLevel() logger.Level { return gv }

func (l *loggerImpl) WithCaller(depth int) logger.Logger {
	var log = *l
	log.callerDepth += depth
	return &log
}

func (l *loggerImpl) WithName(name string) logger.Logger {
	if name == "" {
		return l
	}

	var log = *l
	if log.name == "" {
		log.name = name
	} else {
		log.name = fmt.Sprintf("%s.%s", l.name, name)
	}

	return &log
}

func (l *loggerImpl) WithFields(fields ...logger.Field) logger.Logger {
	var log = *l
	log.fields = append(log.fields, fields...)
	return &log
}

func (l *loggerImpl) WithHooks(hooks ...logger.Hook) logger.Logger {
	var log = *l
	log.hooks = append(log.hooks, hooks...)
	return &log
}

func (l *loggerImpl) Enabled(level logger.Level) bool {
	return gv.Enabled(level)
}

func (l *loggerImpl) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, logCtx{}, l)
}

func (l *loggerImpl) Warn() logger.Entry {
	if !l.Enabled(logger.WARNING) {
		return nil
	}

	return newEntry(l, logger.WARNING)
}

func (l *loggerImpl) Err(err error) logger.Entry {
	if !l.Enabled(logger.ERROR) {
		return nil
	}

	return newEntry(l, logger.ERROR)
}

func (l *loggerImpl) Panic() logger.Entry {
	if !l.Enabled(logger.CRITICAL) {
		return nil
	}

	return newEntry(l, logger.CRITICAL)
}

func (l *loggerImpl) Fatal() logger.Entry {
	if !l.Enabled(logger.CRITICAL) {
		return nil
	}

	return newEntry(l, logger.CRITICAL)
}

func (l *loggerImpl) Info() logger.Entry {
	if !l.Enabled(logger.INFO) {
		return nil
	}

	return newEntry(l, logger.INFO)
}

func (l *loggerImpl) Error() logger.Entry {
	if !l.Enabled(logger.ERROR) {
		return nil
	}

	return newEntry(l, logger.ERROR)
}

func (l *loggerImpl) Print(msg string) {
	if !l.Enabled(logger.DEBUG) {
		return
	}
}

func (l *loggerImpl) Printf(format string, args ...any) {
	if !l.Enabled(logger.DEBUG) {
		return
	}
}
