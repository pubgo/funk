package logutil

import (
	"context"

	"github.com/rs/zerolog"
)

const (
	ModuleName = "module"
	LoggerName = "logger"
)

func Record(evt *zerolog.Event, events ...func(e *zerolog.Event)) *zerolog.Event {
	for _, fn := range events {
		fn(evt)
	}
	return evt
}

func RecordCtx(ctx context.Context, evt *zerolog.Event, events ...func(e *zerolog.Event)) *zerolog.Event {
	for _, fn := range events {
		fn(evt)
	}
	return evt.Ctx(ctx)
}
