package errcheck

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/k0kubun/pp/v3"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/log/logfields"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/prototext"
)

func logErr(ctx context.Context, err error, events ...func(e *zerolog.Event)) {
	if err == nil {
		return
	}

	log.Error(ctx).
		Func(func(e *zerolog.Event) {
			e.Str(logfields.Module, "errcheck")
			e.Str("stack", string(debug.Stack()))
			e.Str(logfields.ErrorDetail, pretty().Sprint(err))
			e.Str(logfields.ErrorID, errors.GetErrorId(err))

			for _, fn := range events {
				fn(e)
			}
		}).
		Str(logfields.Error, err.Error()).
		CallerSkipFrame(2).
		Msgf("%s\n%s\n", err.Error(), prototext.Format(errors.ParseErrToPb(err)))
}

var pretty = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(5)
	return printer
})
