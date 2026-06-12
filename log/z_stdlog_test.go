package log_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/log"
)

func TestStdLog(t *testing.T) {
	log.NewStd(log.GetLogger("with_event").
		WithFields(log.Fields{"hello": "world", "int": 100})).
		Print("dddd")
}

func TestStdLogNilLoggerFallsBackSafely(t *testing.T) {
	assert.NotPanics(t, func() {
		log.NewStd(nil).Print("fallback")
	})
}

func TestStdLogPrintlnAddsSpaces(t *testing.T) {
	var buf bytes.Buffer
	std := log.NewStd(log.Output(&buf))

	std.Println("hello", "world")
	assert.Contains(t, buf.String(), "hello world")
	assert.NotContains(t, buf.String(), "helloworld")
}
