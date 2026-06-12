package log_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/log/slogutil"
)

func TestSlog(t *testing.T) {
	slog.SetDefault(slog.New(log.NewSlog(log.GetLogger(""))))
	slog.Info("ok")
	slog.Info("ok", slogutil.Func(func(evt *zerolog.Event) {
		evt.Str("record", "ok")
	}))
}

func TestSlogNilLoggerFallsBackSafely(t *testing.T) {
	var ctx context.Context

	assert.NotPanics(t, func() {
		h := log.NewSlog(nil)
		assert.True(t, h.Enabled(ctx, slog.LevelInfo))
	})
}

func TestSlogEnabledUsesContextLoggerWithoutPanic(t *testing.T) {
	var buf bytes.Buffer
	var ctx context.Context
	ctx = log.WithLogger(ctx, log.Output(&buf))
	h := log.NewSlog(log.FromCtx(ctx))

	assert.True(t, h.Enabled(ctx, slog.LevelInfo))
	assert.NotPanics(t, func() {
		assert.NoError(t, h.Handle(ctx, slog.NewRecord(testTime(), slog.LevelInfo, "ctx slog", 0)))
	})
	assert.Contains(t, buf.String(), "ctx slog")
}

func testTime() (t time.Time) { return }
