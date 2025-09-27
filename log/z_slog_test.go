package log_test

import (
	"log/slog"
	"testing"

	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/log/slogutil"
	"github.com/rs/zerolog"
)

func TestSlog(t *testing.T) {
	slog.SetDefault(slog.New(log.NewSlog(log.GetLogger(""))))
	slog.Info("ok")
	slog.Info("ok", slogutil.Func(func(evt *zerolog.Event) {
		evt.Str("record", "ok")
	}))
}
