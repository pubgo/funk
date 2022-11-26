package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(cfg Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	writer := cfg.Writer
	if writer == nil {
		writer = os.Stdout
	}

	logger := zerolog.New(writer).Level(level).With().Timestamp().Logger()
	if !cfg.AsJson {
		logger = logger.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = writer
		}))
	}
	log.Logger = logger
	zerolog.DefaultContextLogger = &logger
	return logger
}

func IfDebugLevel(l *zerolog.Logger, fn func(e *zerolog.Event)) {
	if l.GetLevel() != zerolog.DebugLevel {
		return
	}

	fn(l.Debug())
}

func Module(l *zerolog.Logger, module string) zerolog.Logger {
	return l.With().Str("module", module).Logger()
}
