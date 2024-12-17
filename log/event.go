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

//go:linkname putEvent github.com/rs/zerolog.putEvent
func putEvent(e *Event)

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

func mergeEvent(to *Event, from ...*Event) *Event {
	if len(from) == 0 {
		return to
	}

	if to == nil {
		to = zerolog.Dict()
	}

	to1 := convertEvent(to)
	to1.buf = bytes.TrimSpace(bytes.Trim(to1.buf, ","))
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

		if len(to1.buf) == 0 {
			to1.buf = append(to1.buf, '{')
			to1.buf = append(to1.buf, buf...)
		} else {
			to1.buf = append(to1.buf, ","...)
			to1.buf = append(to1.buf, buf...)
		}
	}
	to1.buf = bytes.TrimSpace(to1.buf)
	return to
}
