package log

import (
	"github.com/rs/zerolog"
)

// event 和<zerolog.Event>内存对齐
type event struct {
	buf       []byte
	w         zerolog.LevelWriter
	level     zerolog.Level
	done      func(msg string)
	stack     bool
	ch        []Hook
	skipFrame int
}
