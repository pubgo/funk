package log

import (
	"context"
	
	"github.com/samber/lo"
)

func WithNotice() func(e *Event) {
	return func(e *Event) {
		e.Str("alert", "notice").Bool("critical", true)
	}
}

func RecordErr(ctx context.Context, logs ...Logger) func(err error) error {
	return func(err error) error {
		ctx = lo.If(ctx != nil, ctx).ElseF(context.Background)

		var logger = stdLog
		if len(logs) > 0 {
			logger = logs[0]
		}
		logger.WithCallerSkip(2).Err(err, ctx).Msg("record error log")
		return err
	}
}
