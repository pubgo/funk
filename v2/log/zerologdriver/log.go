package log

import (
	"github.com/pubgo/funk/v2/log"
	"github.com/rs/zerolog"
)

func New(log *zerolog.Logger) log.Logger {
	return &loggerImpl{
		log: log,
	}
}
