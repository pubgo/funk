package slogutil

import (
	"log/slog"

	"github.com/rs/zerolog"
)

type LogFunc func(evt *zerolog.Event)

func Func(fn LogFunc) slog.Attr {
	return slog.Any("func", fn)
}
