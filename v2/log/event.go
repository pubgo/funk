package log

import (
	"bytes"
	"unsafe"

	"github.com/rs/zerolog"
)

// event 和 <zerolog.Event> 内存对齐
type event struct {
	buf []byte
}

//go:linkname putEvent github.com/rs/zerolog.putEvent
func putEvent(e *Event)

func WithEvent(evt *Event) func(e *Event) {
	return func(e *Event) {
		defer putEvent(evt)
		evt1 := convertEvent(evt)
		if len(evt1.buf) == 0 {
			return
		}

		evt1.buf = bytes.TrimLeft(evt1.buf, "{")
		evt1.buf = bytes.TrimRight(evt1.buf, ",")
		if len(evt1.buf) == 0 {
			return
		}

		e1 := convertEvent(e)
		if len(e1.buf) > 0 && len(evt1.buf) > 0 {
			e1.buf = append(e1.buf, ',')
		}
		e1.buf = append(e1.buf, evt1.buf...)
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

func mergeEvent(to *Event, from ...*Event) *Event {
	if len(from) == 0 {
		return to
	}

	if to == nil {
		to = zerolog.Dict()
	}

	to1 := convertEvent(to)
	to1.buf = bytes.TrimRight(to1.buf, ",")
	for i := range from {
		if from[i] == nil {
			continue
		}

		from1 := convertEvent(from[i])
		from1.buf = bytes.TrimLeft(from1.buf, "{")
		from1.buf = bytes.Trim(from1.buf, ",")
		if len(from1.buf) > 0 {
			to1.buf = append(to1.buf, ',')
			to1.buf = append(to1.buf, from1.buf...)
		}
	}
	return to
}
