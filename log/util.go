package log

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/pubgo/funk/v2/internal/errors/errinter"
)

func errDetail(err error) string {
	if err == nil {
		return ""
	}

	return prototext.Format(errinter.ParseErrToPb(err))
}

func RecordErr(logs ...Logger) func(ctx context.Context, err error) error {
	return func(ctx context.Context, err error) error {
		ctx = lo.If(ctx != nil, ctx).ElseF(context.Background)

		logger := stdLog
		if len(logs) > 0 {
			logger = logs[0]
		}
		logger.WithCallerSkip(3).Err(err, ctx).Msg(err.Error())
		return err
	}
}
