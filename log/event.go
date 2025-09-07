package log

import (
	"bytes"
	"slices"
	"unsafe"

	"github.com/rs/zerolog"
)

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

func cloneEvent(target *Event) *Event {
	newTarget := zerolog.Dict()
	convertEvent(newTarget).buf = bytes.Clone(convertEvent(target).buf)
	return newTarget
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
