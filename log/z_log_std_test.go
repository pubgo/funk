package log_test

import (
	"testing"

	"github.com/pubgo/funk/log"
)

func TestStdLog(t *testing.T) {
	evt := log.NewEvent().Str("hello", "world").Int("int", 100)
	log.NewStd(log.GetLogger("with_event").WithEvent(evt)).Print("dddd")
}
