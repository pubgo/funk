package log

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
	slogcommon "github.com/samber/slog-common"
)

func NewSlog(log Logger) slog.Handler {
	return &slogImpl{l: log.WithCallerSkip(3)}
}

var logLevels = map[slog.Level]zerolog.Level{
	slog.LevelDebug: zerolog.DebugLevel,
	slog.LevelInfo:  zerolog.InfoLevel,
	slog.LevelWarn:  zerolog.WarnLevel,
	slog.LevelError: zerolog.ErrorLevel,
}
var _ slog.Handler = (*slogImpl)(nil)

type slogImpl struct {
	l Logger
}

func (s slogImpl) Enabled(ctx context.Context, level slog.Level) bool {
	return s.l.(*loggerImpl).enabled(ctx, logLevels[level])
}

func (s slogImpl) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < 0 {
		r.Level = slog.LevelDebug
	}

	logger := s.l.WithLevel(logLevels[r.Level])

	var evt *Event
	switch r.Level {
	case slog.LevelDebug:
		evt = logger.Debug(ctx)
	case slog.LevelInfo:
		evt = logger.Info(ctx)
	case slog.LevelWarn:
		evt = logger.Warn(ctx)
	case slog.LevelError:
		evt = logger.Error(ctx)
	}

	if evt == nil {
		return nil
	}

	if !r.Time.IsZero() {
		evt.Time(zerolog.TimestampFieldName, r.Time)
	}

	r.Attrs(func(attr slog.Attr) bool {
		evt.Any(attr.Key, attr.Value.Any())
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
