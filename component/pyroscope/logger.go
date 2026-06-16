package pyroscope

import (
	"github.com/pubgo/funk/v2/log"
)

type loggerAdapter struct {
	logger log.Logger
}

func newLoggerAdapter(logger log.Logger) *loggerAdapter {
	return &loggerAdapter{logger: logger.WithName("pyroscope")}
}

func (l *loggerAdapter) Infof(format string, args ...any) {
	l.logger.Info().Msgf(format, args...)
}

func (l *loggerAdapter) Debugf(format string, args ...any) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *loggerAdapter) Errorf(format string, args ...any) {
	l.logger.Error().Msgf(format, args...)
}
