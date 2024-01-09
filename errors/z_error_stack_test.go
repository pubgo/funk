package errors

import (
	"testing"

	"github.com/pubgo/funk/stack"
)

func TestStack(t *testing.T) {
	var err = WrapCaller(New("hello error"))
	err = Wrap(err, "next error")
	err = WrapTag(err,
		T("event", "test event"),
		T("test123", 123),
		T("test", "hello"),
		T("fn_stack", stack.CallerWithFunc(stack.CallerWithFunc)),
	)

	err = WrapStack(err)
	err = Wrapf(err, "next error name=%s", "wrapf")
	Debug(err)
}
