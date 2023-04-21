package errors

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	var err = WrapCaller(New("hello error"))
	err = Wrap(err, "next error")
	//err = WrapEventFn(err, func(evt *Event) {
	//	evt.Str("event", "test event")
	//	evt.Int64("test123", 123)
	//	evt.Str("test", "hello")
	//})
	err = WrapStack(err)
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = Append(err, fmt.Errorf("raw error"))
	err = Append(err, New("New errors error"))
	err = Append(err, SimpleErr(func(err *Err) {
		err.Err = fmt.Errorf("test Err")
		err.Msg = "Err errors error"
		err.Tags = map[string]any{"tags": "hello"}
	}))
	Debug(err)
}
