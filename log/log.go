package log

import (
	"github.com/pubgo/funk/version"
	"github.com/rs/zerolog"
)

func NewStd(log Logger) StdLogger {
	return &stdLogImpl{log: log.WithCallerSkip(1)}
}

func New(log *zerolog.Logger) Logger {
	return &loggerImpl{
		log:  log,
		name: version.Project(),
	}
}
