package log

// event 和<zerolog.Event>内存对齐
type event struct {
	buf []byte
}
