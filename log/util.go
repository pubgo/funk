package log

import (
	"bytes"
	"context"
	"slices"
	"unsafe"

	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/debugs"
	"github.com/pubgo/funk/v2/errors"
)

func errDetail(err error) string {
	if err == nil {
		return ""
	}

	if debugs.Enabled.Value() {
		errors.DebugPrint(err)
	}

	return string(errors.JsonPrint(err))
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

// event 和 <zerolog.Event> 内存对齐
type event struct {
	buf []byte
}

func WithEvent(evt *Event) func(e *Event) {
	return func(e *Event) {
		if !e.Enabled() {
			return
		}

		buf := slices.Clone(convertEvent(evt).buf)
		if len(buf) == 0 {
			return
		}

		buf = bytes.TrimLeft(buf, "{")
		buf = bytes.TrimSpace(bytes.Trim(buf, ","))
		if len(buf) == 0 {
			return
		}

		e1 := convertEvent(e)
		e1.buf = bytes.TrimSpace(e1.buf)
		if len(e1.buf) == 0 {
			e1.buf = append(e1.buf, '{')
			e1.buf = append(e1.buf, buf...)
		} else {
			e1.buf = append(e1.buf, ","...)
			e1.buf = append(e1.buf, buf...)
		}

		e1.buf = bytes.TrimSpace(e1.buf)
	}
}

func convertEvent(event2 *Event) *event {
	return (*event)(unsafe.Pointer(event2))
}

func NewEvent() *Event {
	return zerolog.Dict()
}

func GetEventBuf(evt *Event) []byte {
	if evt == nil {
		return nil
	}

	return append(convertEvent(evt).buf, '}')
}

func mergeEvent(target *Event, from ...*Event) *Event {
	if len(from) == 0 {
		return target
	}

	if target == nil {
		target = zerolog.Dict()
	}

	targetEvent := convertEvent(target)
	targetEvent.buf = bytes.TrimSpace(bytes.Trim(targetEvent.buf, ","))
	for i := range from {
		if from[i] == nil {
			continue
		}

		buf := slices.Clone(convertEvent(from[i]).buf)
		if len(buf) == 0 {
			continue
		}

		buf = bytes.TrimLeft(buf, "{")
		buf = bytes.TrimSpace(bytes.Trim(buf, ","))
		if len(buf) == 0 {
			continue
		}

		if len(targetEvent.buf) == 0 {
			targetEvent.buf = append(targetEvent.buf, '{')
			targetEvent.buf = append(targetEvent.buf, buf...)
		} else {
			targetEvent.buf = append(targetEvent.buf, ","...)
			targetEvent.buf = append(targetEvent.buf, buf...)
		}
	}
	targetEvent.buf = bytes.TrimSpace(targetEvent.buf)
	return target
}
