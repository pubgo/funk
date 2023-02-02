package log_test

import (
	"errors"
	"testing"

	"github.com/pubgo/funk/log"
)

func TestName(t *testing.T) {
	log.Debug().Str("hello", "world world").Msg("ok ok")
	log.Print("world world")
	log.Info().Str("hello", "world world").Msg("ok ok")
	log.Warn().Str("hello", "world world").Msg("ok ok")
	log.Err(errors.New("error")).Str("hello", "world world").Msg("ok ok")
	log.GetLogger("test_app").Info().Str("hello", "world world").Msg("ok ok")
	log.GetLogger("test_app").Info().Str("hello", "world world").Msg("ok ok")
	log.GetLogger("test_app").
		WithFields(log.Map{"module": "pkg"}).
		Info().
		Str("hello", "world world").
		Func(log.WithNotice()).
		Msg("ok ok")
}

func TestEvent(t *testing.T) {
	var evt = log.NewEvent().Str("hello", "world").Int("int", 100)
	ee := log.Info().Str("info", "abcd")
	ee.Func(log.WithEvent(evt))
	ee.Msg("dddd")
}

func TestWithEvent(t *testing.T) {
	var evt = log.NewEvent().Str("hello", "world").Int("int", 100)
	ee := log.GetLogger("with_event").WithEvent(evt).Info().Str("info", "abcd")
	ee.Msg("dddd")
}
