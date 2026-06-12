package log

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
	slogcommon "github.com/samber/slog-common"

	"github.com/pubgo/funk/v2/log/slogutil"
)

func NewSlog(log Logger) slog.Handler {
	if log == nil {
		log = stdLog
	}

	return &slogImpl{l: log.WithCallerSkip(3)}
}

var logLevels = map[slog.Level]zerolog.Level{
	slog.LevelDebug: zerolog.DebugLevel,
	slog.LevelInfo:  zerolog.InfoLevel,
	slog.LevelWarn:  zerolog.WarnLevel,
	slog.LevelError: zerolog.ErrorLevel,
}

func convertSlog(lvl slog.Level) slog.Level {
	switch {
	case lvl < slog.LevelInfo:
		return slog.LevelDebug
	case lvl < slog.LevelWarn:
		return slog.LevelInfo
	case lvl < slog.LevelError:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

var _ slog.Handler = (*slogImpl)(nil)

type slogImpl struct {
	l Logger
}

func (s slogImpl) Enabled(ctx context.Context, level slog.Level) bool {
	if isLogDisabled(ctx) {
		return false
	}

	if enabler, ok := s.l.(interface {
		enabled(context.Context, zerolog.Level) bool
	}); ok {
		return enabler.enabled(ctx, logLevels[convertSlog(level)])
	}

	return logLevels[convertSlog(level)] >= zerolog.GlobalLevel()
}

func (s slogImpl) Handle(ctx context.Context, r slog.Record) error {
	if isLogDisabled(ctx) {
		return nil
	}

	logger := s.l
	level := convertSlog(r.Level)
	var evt *Event
	switch level {
	case slog.LevelDebug:
		evt = logger.Debug(ctx)
	case slog.LevelInfo:
		evt = logger.Info(ctx)
	case slog.LevelWarn:
		evt = logger.Warn(ctx)
	default:
		evt = logger.Error(ctx)
	}
	if evt == nil {
		return nil
	}

	if !r.Time.IsZero() {
		evt.Time(zerolog.TimestampFieldName, r.Time)
	}

	r.Attrs(func(attr slog.Attr) bool {
		if fn, ok := attr.Value.Any().(slogutil.LogFunc); ok {
			evt.Func(fn)
		} else {
			evt.Any(attr.Key, attr.Value.Any())
		}
		return true
	})

	evt.Msg(r.Message)
	return nil
}

func (s slogImpl) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogImpl{l: s.l.WithFields(slogcommon.AttrsToMap(attrs...))}
}

func (s slogImpl) WithGroup(name string) slog.Handler {
	return &slogImpl{l: s.l.WithName(name)}
}
