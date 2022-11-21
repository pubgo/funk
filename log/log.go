package log

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pubgo/funk/logger"
)

type Logger struct {
	callDepth int
	level     uint
	name      string
	tags      []logger.Tagger
}

func (l Logger) Enabled() bool {
	return l.level <= ll
}

func (l Logger) V(level uint) Logger {
	if level == 0 {
		return l
	}

	l.level = level
	return l
}

func (l Logger) WithName(name string) Logger {
	if name == "" {
		return l
	}

	if l.name == "" {
		l.name = name
	} else {
		l.name = l.name + NameDelim + name
	}

	return l
}

func (l Logger) WithTags(tags ...logger.Tagger) Logger {
	if len(tags) == 0 {
		return l
	}

	l.tags = append(l.tags, tags...)
	return l
}

func (l Logger) WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, logCtx{}, l)
}

func (l Logger) Info(msg string, tags ...logger.Tagger) {
	if !l.Enabled() {
		return
	}

	tags = append(tags, logger.Tag("logger", l.name))
	tags = append(tags, logger.Tag("level", "info"))
	tags = append(tags, logger.Tag("msg", msg))
	tags = append(tags, logger.Tag("ts", time.Now().UTC().String()))

	for i := range hooks {
		tags = hooks[i].Hook(tags)
	}

	writer.Info(l.level, msg, tags...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	if !l.Enabled() {
		return
	}

	writer.Info(l.level, fmt.Sprintf(format, args...))
}

func (l Logger) Error(err error, msg string, tags ...logger.Tagger) {
	if err == nil || reflect.ValueOf(err).IsZero() {
		return
	}

	for i := range hooks {
		tags = hooks[i].Hook(tags)
	}

	writer.Error(err, msg, tags...)
}
