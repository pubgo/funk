package log

import (
	"context"
)

func WithNotice() func(e *Event) {
	return func(e *Event) {
		e.Str("alert", "notice").Bool("critical", true)
	}
}

func RecordErr(logs ...Logger) func(ctx context.Context, err error) error {
	return func(ctx context.Context, err error) error {
		var logger = stdLog
		if len(logs) > 0 {
			logger = logs[0]
		}
		logger.WithCallerSkip(2).Err(err, ctx).Msg("record error log")
		return err
	}
}
