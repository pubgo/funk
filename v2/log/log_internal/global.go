package log_internal

import (
	"context"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/generic"
)

var (
	logEnableChecker      = func(ctx context.Context, lvl Level, nameOrMessage string, fields Map) bool { return true }
	zErrMarshalFunc       = zerolog.ErrorMarshalFunc
	zInterfaceMarshalFunc = zerolog.InterfaceMarshalFunc
	logGlobalHook         = zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		if logEnableChecker == nil {
			return
		}

		if logEnableChecker(e.GetCtx(), level, message, nil) {
			return
		}

		e.Discard()
	})
	_ = generic.Init(func() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.ErrorMarshalFunc = func(err error) interface{} {
			if err == nil {
				return nil
			}

			switch err.(type) {
			case json.Marshaler:
				return &logLogObjectMarshaler{err: err}
			}

			if zErrMarshalFunc == nil {
				return err.Error()
			}

			return zErrMarshalFunc(err)
		}

		zerolog.InterfaceMarshalFunc = func(v any) ([]byte, error) {
			if v == nil {
				return []byte("null"), nil
			}

			switch e1 := v.(type) {
			case json.Marshaler:
				return e1.MarshalJSON()
			}

			return zInterfaceMarshalFunc(v)
		}
	})

	// stdZeroLog default zerolog for debug
	stdZeroLog = generic.Ptr(
		zerolog.New(os.Stderr).
			Level(zerolog.DebugLevel).
			With().Timestamp().
			Caller().Logger().
			Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
				w.Out = os.Stderr
				w.TimeFormat = time.RFC3339
			})).Hook(new(hookImpl), logGlobalHook),
	)

	_ = generic.Init(func() {
		zlog.Logger = generic.FromPtr(stdZeroLog)
	})

	// stdLog is the global logger.
	stdLog = New(nil)
)

func GetLogger() Logger { return stdLog }

// SetLogger set global log
func SetLogger(log *zerolog.Logger) {
	assert.If(log == nil, "[log] should not be nil")

	log = generic.Ptr(log.Hook(logGlobalHook))

	stdZeroLog = log
	zlog.Logger = *log
}

func SetEnableChecker(checker EnableChecker) {
	if checker == nil {
		return
	}

	logEnableChecker = checker
}
