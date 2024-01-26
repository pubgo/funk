package log_test

import (
	"context"
	"fmt"
	"github.com/pubgo/funk/errors"
	"testing"

	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
)

func TestName(t *testing.T) {
	log.Debug().Str("hello", "world world").Msg("ok ok")
	log.Print("world world")
	log.Info().Str("hello", "world world").Msg("ok ok")
	log.Warn().Str("hello", "world world").Msg("ok ok")

	var err = errors.WrapCaller(fmt.Errorf("test error"))
	err = errors.Wrap(err, "next error")
	err = errors.WrapTag(err, errors.T("event", "test event"), errors.T("test123", 123), errors.T("test", "hello"))
	err = errors.Wrapf(err, "next error name=%s", "wrapf")
	log.Err(err).Str("hello", "world world").Msg("ok ok")
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
	var evt = log.NewEvent().Str("hello", "world").Int("int", 100).Dict("ddd", log.NewEvent())
	ctx := log.CreateEventCtx(context.Background(), evt)
	ee := log.Info(ctx).Str("info", "abcd").Func(log.WithEvent(evt))
	ee.Msg("dddd")
}

func TestWithEvent(t *testing.T) {
	var evt = log.NewEvent().Str("hello", "hello world").Int("int", 100)
	ee := log.GetLogger("with_event").WithEvent(evt).Info().Str("info", "abcd")
	ee.Msg("dddd")
}

func TestSetLog(t *testing.T) {
	//zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.SetLogger(generic.Ptr(zl.Output(zerolog.NewConsoleWriter())))
	log.Debug().Msg("test")
	log.Info().Msg("test")
	log.Warn().Msg("test")
	log.Error().Msg("test")
}
