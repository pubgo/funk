package log

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logCtxString string

var logCtx = logCtxString("log:" + time.Now().String())

func CreateCtx(ctx context.Context, log zerolog.Logger) context.Context {
	return context.WithValue(ctx, logCtx, log)
}

func GetLog(ctx context.Context) *zerolog.Logger {
	if l, ok := ctx.Value(logCtx).(zerolog.Logger); ok {
		return &l
	}

	return &log.Logger
}
