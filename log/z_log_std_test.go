package log_test

import (
	"testing"

	"github.com/pubgo/funk/v2/log"
)

func TestStdLog(t *testing.T) {
	log.NewStd(log.GetLogger("with_event").
		WithFields(log.Map{"hello": "world", "int": 100})).
		Print("dddd")
}
