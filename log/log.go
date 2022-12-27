package log

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pubgo/funk/logger"
)

// A Logger represents an active logging object that generates lines of JSON output to an io.Writer.
type Logger1 struct {
	// Level defines log levels.
	Level logger.Level

	// Caller determines if adds the file:line of the "caller" key.
	// If Caller is negative, adds the full /path/to/file:line of the "caller" key.
	Caller int

	// TimeField defines the time filed name in output.  It uses "time" in if empty.
	TimeField string

	// TimeFormat specifies the time format in output. It uses time.RFC3339 with milliseconds if empty.
	// If set with `TimeFormatUnix`, `TimeFormatUnixMs`, times are formated as UNIX timestamp.
	TimeFormat string

	// Context specifies an optional context of logger.
	Context Context

	// Writer specifies the writer of output. It uses a wrapped os.Stderr Writer in if empty.
	Writer Writer
}

type Logger struct {
	callDepth int
	name      string
	tags      []logger.Tagger
}

func (l Logger) Enabled(level logger.Level) bool {
	return level <= gv
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
	if !l.Enabled(logger.INFO) {
		return
	}

	tags = append(tags, logger.Tag("logger", l.name))
	tags = append(tags, logger.Tag("level", "info"))
	tags = append(tags, logger.Tag("msg", msg))
	tags = append(tags, logger.Tag("ts", time.Now().UTC().String()))

	for i := range hooks {
		tags = hooks[i].Hook(tags)
	}

	writer.Log(tags...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	if !l.Enabled(logger.INFO) {
		return
	}

	writer.Log(fmt.Sprintf(format, args...))
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
