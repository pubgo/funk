package logutil

import (
	"context"

	"github.com/rs/zerolog"
)

// Record adds a notice alert to the event.
func Record(evt *zerolog.Event, events ...func(e *zerolog.Event)) *zerolog.Event {
	for _, fn := range events {
		fn(evt)
	}
	return evt
}

// RecordCtx adds a notice alert to the event.
func RecordCtx(ctx context.Context, evt *zerolog.Event, events ...func(e *zerolog.Event)) *zerolog.Event {
	for _, fn := range events {
		fn(evt)
	}
	return evt.Ctx(ctx)
}

// WithNotice adds a notice alert to the event.
func WithNotice() func(e *zerolog.Event) {
	return func(e *zerolog.Event) {
		e.Str("alert", "notice").Bool("critical", true)
	}
}
