package log

import _ "unsafe"

//go:linkname putEvent github.com/rs/zerolog.putEvent
func putEvent(e *Event)
