package log

import (
	"sync/atomic"

	"github.com/phuslu/goid"
	"github.com/rs/zerolog"
)

var _ zerolog.Hook = (*hookImpl)(nil)

type hookImpl struct {
	count uint64
}

func (h *hookImpl) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if zerolog.GlobalLevel() >= zerolog.WarnLevel {
		return
	}

	e.Int64("go_id", goid.Goid())
	e.Uint64("log_count", atomic.AddUint64(&h.count, 1))
}
