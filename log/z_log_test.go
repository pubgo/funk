package log_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/log/logutil"
)

func TestLogLevel(t *testing.T) {
	log.Warn().Msg("test warn")
	logger := log.GetLogger().WithLevel(zerolog.ErrorLevel)
	logger.Warn().Msg("test warn")
}

func TestWithName(t *testing.T) {
	log.GetLogger("log1").
		Debug().
		Func(func(e *zerolog.Event) {
			buf := gjson.ParseBytes(log.GetEventBuf(e))
			assert.Equal(t, buf.Get("logger").String(), "log1")
		}).Msg("hello")

	log.GetLogger("log1").
		WithName("log2").
		Debug().
		Func(func(e *zerolog.Event) {
			buf := gjson.ParseBytes(log.GetEventBuf(e))
			assert.Equal(t, buf.Get("logger").String(), "log1.log2")
		}).Msg("hello")

	log.Debug().
		Func(func(e *zerolog.Event) {
			buf := gjson.ParseBytes(log.GetEventBuf(e))
			assert.Equal(t, buf.Get("logger").String(), "")
		}).Msg("hello")
}

func TestNilLog(t *testing.T) {
	var buf bytes.Buffer
	log.Output(&buf).Debug().Any("key", nil).Send()
	ret := gjson.ParseBytes(buf.Bytes())
	assert.Equal(t, ret.Get("key").String(), "")

	log.OutputWriter(func(p []byte) (n int, err error) {
		parseBytes := gjson.ParseBytes(buf.Bytes())
		assert.Equal(t, parseBytes.Get("key").String(), "")
		return len(p), nil
	}).Debug().Any("key", nil).Send()
}

func TestWithDisabled(t *testing.T) {
	ctx := log.WithDisabled(context.Background())
	evt := log.Info(ctx).Str("hello", "world world")
	assert.Equal(t, string(log.GetEventBuf(evt)), "")
}

func TestName(t *testing.T) {
	log.Debug().Str("hello", "world world").Msg("ok ok")
	log.Print("world world")
	log.Info().Str("hello", "world world").Msg("ok ok")
	log.Warn().Str("hello", "world world").Msg("ok ok")

	err := errors.WrapCaller(fmt.Errorf("test error"))
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
		Func(logutil.WithNotice()).
		Msg("ok ok")
}

func TestEvent(t *testing.T) {
	getEvt := func() *log.Event {
		return log.NewEvent().Str("hello", "world").Int("int", 100).Dict("ddd", log.NewEvent())
	}

	getCtx := func(evt *log.Event) context.Context {
		return log.CreateEventCtx(context.Background(), evt)
	}

	t.Run("event ctx", func(t *testing.T) {
		log.Info(getCtx(getEvt())).Send()
	})

	t.Run("event func", func(t *testing.T) {
		log.Info().Func(log.WithEvent(getEvt())).Send()
	})

	t.Run("update event ctx", func(t *testing.T) {
		log.Info(log.UpdateEventCtx(getCtx(getEvt()), log.Map{"add-update-event": "ok"})).Send()
	})
}

func TestWithEvent(t *testing.T) {
	ee := log.GetLogger("with_event").
		WithFields(log.Map{"hello": "hello world", "int": 100}).
		Info().
		Str("info", "abcd")
	ee.Msg("dddd")
}

func TestSetLog(t *testing.T) {
	logger := log.Output(zerolog.NewConsoleWriter())
	logger.Debug().Msg("test")
	logger.Info().Msg("test")
	logger.Warn().Msg("test")
	logger.Error().Msg("test")
}

func TestChecker(t *testing.T) {
	l := log.GetLogger("test-checker")
	l.Info().Msg("hello")

	log.SetEnableChecker(func(ctx context.Context, lvl log.Level, name, message string, fields log.Map) bool {
		fmt.Println(lvl, name, fields)
		return true
	})
	l.Info().Msg("hello1")
	l.Warn().Msg("hello1")
	l.Error().Msg("hello1")
	l.Debug().Msg("hello1")
}

func TestErr(t *testing.T) {
	err := fmt.Errorf("test error")
	log.Error().Err(err).Msg(err.Error())

	err1 := errors.Errorf("test format")
	log.Error().Err(err1).Msg(err1.Error())
}

func TestError(t *testing.T) {
	err := fmt.Errorf("test raw error")
	log.Error().Err(err).Msg(err.Error())

	err1 := errors.Errorf("test errors format")
	log.Error().Err(err1).Msg("raw error: " + err1.Error())
	log.Err(err1).Msg(err1.Error())
}
