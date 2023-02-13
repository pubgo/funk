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
	var err = WrapCaller(New("hello error"))
	err = Wrap(err, "next error")
	err = WrapEventFn(err, func(evt *Event) {
		evt.Str("event", "test event")
		evt.Int64("test123", 123)
		evt.Str("test", "hello")
	})
	err = Wrapf(err, "next error name=%s", "wrapf")
	err = WrapCode(err, errorpb.Code_Canceled)
	err = WrapReason(err, "test reason")
	err = Append(err, fmt.Errorf("raw error"))
	err = Append(err, New("New errors error"))
	err = Append(err, SimpleErr(func(err *Err) {
		err.Err = fmt.Errorf("test Err")
		err.Msg = "Err errors error"
		err.Tags = map[string]any{"tags": "hello"}
	}))

	err = WrapBizCode(err, "user.not_found")
	err = WrapCodeFn(err, func(err ErrCode) {
		err.SetReason("user not_found")
	})

	err = WrapStack(err)
	IfErr(err, func(err error) {
		Debug(err)
	})

	var fff ErrCode
	t.Log(As(err, &fff))
	t.Log(fff.Reason())
}
