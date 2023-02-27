package errors

import (
	"fmt"
	"testing"

	"github.com/pubgo/funk/proto/errorpb"
)

func init1() error {
	return WrapCaller(New("test skip"), 1)
}

func TestSkip(t *testing.T) {
	var err = init1()
	Debug(err)
}

func TestFormat(t *testing.T) {
	var err = WrapCaller(fmt.Errorf("test error, err=%w", New("hello error")))
	err = Wrap(err, "next error")
	err = WrapEventFn(err, func(evt *Event) {
		evt.Str("event", "test event")
		evt.Int64("test123", 123)
		evt.Str("test", "hello")
	})
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = Append(err, fmt.Errorf("raw error"))
	err = Append(err, New("New errors error"))
	err = Append(err, SimpleErr(func(err *Err) {
		err.Err = fmt.Errorf("test Err")
		err.Msg = "Err errors error"
		err.Tags = map[string]any{"tags": "hello"}
	}))

	err = NewCode(errorpb.Code_Canceled).
		SetReason("user not_found").
		SetName("user.not_found").
		SetStatus(100).
		SetErr(err).AddTag("hello", "world")

	err = WrapStack(err)
	IfErr(err, func(err error) {
		Debug(err)
	})

	var fff ErrCode
	t.Log(As(err, &fff))
	t.Log(fff.Reason())
}
