package log

import (
	"unsafe"

	"github.com/rs/zerolog"
)

func WithNotice() func(e *Event) {
	return func(e *Event) {
		e.Str("alert", "notice").Bool("critical", true)
	}
}

func WithEvent(evt *Event) func(e *Event) {
	return func(e *Event) {
		evt1 := convertEvent(evt)
		if evt1.buf[0] == '{' && len(evt1.buf) == 1 {
			return
		}

		e1 := convertEvent(e)
		e1.buf = append(e1.buf, ',')
		e1.buf = append(e1.buf, evt1.buf[1:]...)
		putEvent(evt)
	}
}

func convertEvent(event2 *Event) *event {
	return (*event)(unsafe.Pointer(event2))
}

func NewEvent() *Event {
	return zerolog.Dict()
}
