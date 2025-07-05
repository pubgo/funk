package log

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kr/pretty"
	"github.com/samber/lo"
)

func WithNotice() func(e *Event) {
	return func(e *Event) {
		e.Str("alert", "notice").Bool("critical", true)
	}
}

func RecordErr(logs ...Logger) func(ctx context.Context, err error) error {
	return func(ctx context.Context, err error) error {
		ctx = lo.If(ctx != nil, ctx).ElseF(context.Background)

		var logger = stdLog
		if len(logs) > 0 {
			logger = logs[0]
		}
		logger.WithCallerSkip(3).Err(err, ctx).Msg(err.Error())
		return err
	}
}

func errDetail(err error) []byte {
	if err == nil {
		return nil
	}

	switch errData := err.(type) {
	case json.Marshaler:
		data, err1 := errData.MarshalJSON()
		if err1 != nil {
			return []byte(fmt.Sprintf("%s: %s", err1.Error(), pretty.Sprint(err)))
		}
		return data
	default:
		return []byte(pretty.Sprint(err))
	}
}
