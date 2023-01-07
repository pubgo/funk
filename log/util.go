package log

func WithNotice() func(e *Event) {
	return func(e *Event) {
		e.Str("alert", "notice")
	}
}
