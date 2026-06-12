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

func TestWithLoggerHandlesNilInputs(t *testing.T) {
	var buf bytes.Buffer
	var ctx context.Context
	ctx = log.WithLogger(ctx, log.Output(&buf))

	log.FromCtx(ctx).Info().Msg("with logger helper")
	assert.Contains(t, buf.String(), "with logger helper")

	buf.Reset()
	log.Info(ctx).Msg("global helper uses ctx logger")
	assert.Contains(t, buf.String(), "global helper uses ctx logger")

	assert.NotPanics(t, func() {
		assert.NotNil(t, log.FromCtx(log.WithLogger(ctx, nil)))
	})
}

func TestFromCtxFallsBackToDefaultLogger(t *testing.T) {
	var ctx context.Context
	assert.NotNil(t, log.FromCtx(ctx))
}

func TestName(t *testing.T) {
	log.Debug().Str("hello", "world world").Msg("ok ok")
	log.Print("world world")
	log.Info().Str("hello", "world world").Msg("ok ok")
	log.Warn().Str("hello", "world world").Msg("ok ok")

	err := errors.WrapCaller(fmt.Errorf("test error"))
	err = errors.Wrap(err, "next error")
	err = errors.WrapTags(err, errors.Tags{"event": "test event", "test123": 123, "test": "hello"})
	err = errors.Wrapf(err, "next error name=%s", "wrapf")
	log.Err(err).Str("hello", "world world").Msg("ok ok")
	log.GetLogger("test_app").Info().Str("hello", "world world").Msg("ok ok")
	log.GetLogger("test_app").Info().Str("hello", "world world").Msg("ok ok")
	log.GetLogger("test_app").
		WithFields(log.Fields{"module": "pkg"}).
		Info().
		Str("hello", "world world").
		Func(logutil.WithNotice()).
		Msg("ok ok")
}

func TestEvent(t *testing.T) {
	getEvt := func() log.Fields {
		return log.Fields{
			"hello": "world",
			"int":   100,
			"float": 1.23,
		}
	}

	getCtx := func(evt log.Fields) context.Context {
		return log.CreateFieldsCtx(context.Background(), evt)
	}

	t.Run("event ctx", func(t *testing.T) {
		log.Info(getCtx(getEvt())).Send()
	})

	t.Run("update event ctx", func(t *testing.T) {
		log.Info(log.UpdateFieldsCtx(getCtx(getEvt()), log.Fields{"add-update-event": "ok"})).Send()
	})

	t.Run("with fields helper", func(t *testing.T) {
		ctx := log.WithFields(getCtx(getEvt()), log.Fields{"via": "helper"})
		log.Info(ctx).Send()
		assert.Equal(t, "helper", log.GetFieldsFromCtx(ctx)["via"])
	})
}

func TestWithFieldsHandlesNilContext(t *testing.T) {
	var ctx context.Context
	ctx = log.WithFields(ctx, log.Fields{"request_id": "req-nil"})
	fields := log.GetFieldsFromCtx(ctx)

	assert.Equal(t, "req-nil", fields["request_id"])
}

func TestGetFieldsFromCtxHandlesNilContext(t *testing.T) {
	var ctx context.Context

	assert.NotPanics(t, func() {
		assert.Nil(t, log.GetFieldsFromCtx(ctx))
	})
}

func TestCreateFieldsCtxClonesInputMap(t *testing.T) {
	source := log.Fields{"request_id": "req-source", "scope": "original"}
	ctx := log.CreateFieldsCtx(context.Background(), source)

	source["request_id"] = "req-mutated"
	source["added_later"] = true

	fields := log.GetFieldsFromCtx(ctx)
	assert.Equal(t, "req-source", fields["request_id"])
	assert.Equal(t, "original", fields["scope"])
	assert.False(t, gjson.ParseBytes(mustEventJSON(t, fields)).Get("added_later").Exists())
}

func TestUpdateFieldsCtxDoesNotMutateParent(t *testing.T) {
	base := log.CreateFieldsCtx(context.Background(), log.Fields{"request_id": "req-1"})
	updated := log.UpdateFieldsCtx(base, log.Fields{"user_id": 42})

	baseFields := log.GetFieldsFromCtx(base)
	updatedFields := log.GetFieldsFromCtx(updated)

	assert.Equal(t, "req-1", baseFields["request_id"])
	assert.Nil(t, baseFields["user_id"])
	assert.Equal(t, "req-1", updatedFields["request_id"])
	assert.Equal(t, float64(42), gjson.ParseBytes(mustEventJSON(t, updatedFields)).Get("user_id").Float())
}

func TestWithFieldsLatestWinsWithoutMutatingBaseLogger(t *testing.T) {
	base := log.GetLogger("fields-base").WithFields(log.Fields{"component": "base", "shared": "base"})
	child := base.WithFields(log.Fields{"shared": "child", "request_id": "req-1"})

	childEvent := eventJSON(t, child.Info())
	assert.Equal(t, "base", childEvent.Get("component").String())
	assert.Equal(t, "child", childEvent.Get("shared").String())
	assert.Equal(t, "req-1", childEvent.Get("request_id").String())

	baseEvent := eventJSON(t, base.Info())
	assert.Equal(t, "base", baseEvent.Get("component").String())
	assert.Equal(t, "base", baseEvent.Get("shared").String())
	assert.False(t, baseEvent.Get("request_id").Exists())
}

func TestContextFieldsDoNotPolluteLoggerDefaults(t *testing.T) {
	logger := log.GetLogger("ctx-merge").WithFields(log.Fields{"component": "worker", "scope": "logger"})
	ctx := log.CreateFieldsCtx(context.Background(), log.Fields{"scope": "ctx", "request_id": "req-2"})

	withCtx := eventJSON(t, logger.Info(ctx))
	assert.Equal(t, "worker", withCtx.Get("component").String())
	assert.Equal(t, "ctx", withCtx.Get("scope").String())
	assert.Equal(t, "req-2", withCtx.Get("request_id").String())

	withoutCtx := eventJSON(t, logger.Info())
	assert.Equal(t, "worker", withoutCtx.Get("component").String())
	assert.Equal(t, "logger", withoutCtx.Get("scope").String())
	assert.False(t, withoutCtx.Get("request_id").Exists())
}

func TestWithEvent(t *testing.T) {
	ee := log.GetLogger("with_event").
		WithFields(log.Fields{"hello": "hello world", "int": 100}).
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
	t.Cleanup(func() {
		log.SetEnableChecker(nil)
	})

	l := log.GetLogger("test-checker")
	l.Info().Msg("hello")

	log.SetEnableChecker(func(ctx context.Context, lvl log.Level, name, message string, fields log.Fields) bool {
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

func TestRecordErrIgnoresNil(t *testing.T) {
	var buf bytes.Buffer
	record := log.RecordErr(log.Output(&buf))

	assert.NotPanics(t, func() {
		assert.Nil(t, record(context.Background(), nil))
	})
	assert.Empty(t, buf.String())
}

func eventJSON(t *testing.T, evt *log.Event) gjson.Result {
	t.Helper()

	return gjson.ParseBytes(log.GetEventBuf(evt))
}

func mustEventJSON(t *testing.T, fields log.Fields) []byte {
	t.Helper()

	evt := log.NewEvent().Fields(fields)
	return log.GetEventBuf(evt)
}
