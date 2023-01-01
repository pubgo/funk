package log

import (
	"context"
	_ "github.com/phuslu/log"
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

func (l *loggerImpl) Print(msg string) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerImpl) Printf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerImpl) GetLevel() logger.Level { return gv }

func (l *loggerImpl) WithCaller(depth int) logger.Logger {
	l.callerDepth += depth
	return l
}

func (l *loggerImpl) WithName(name string) logger.Logger {
	if name == "" {
		return l
	}

	if l.name == "" {
		l.name = name
	} else {
		l.name = l.name + "." + name
	}

	return l
}

func (l *loggerImpl) WithFields(fields ...logger.Field) logger.Logger {
	l.fields = append(l.fields, fields...)
	return l
}

func (l *loggerImpl) WithHooks(hooks ...logger.Hook) logger.Logger {
	l.hooks = append(l.hooks, hooks...)
	return l
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

}

func (l *loggerImpl) Err(err error) logger.Entry {
	if !l.Enabled(logger.ERROR) {
		return nil
	}

}

func (l *loggerImpl) Panic() logger.Entry {
	if !l.Enabled(logger.CRITICAL) {
		return nil
	}

}

func (l *loggerImpl) Fatal() logger.Entry {
	if !l.Enabled(logger.CRITICAL) {
		return nil
	}

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

}
