package log

import (
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
