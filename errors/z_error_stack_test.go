package errors_test

import (
	"testing"

	"github.com/pubgo/funk/errors"

	"github.com/pubgo/funk/stack"
)

func TestStack(t *testing.T) {
	err := errors.WrapCaller(errors.New("hello error"))
	err = errors.Wrap(err, "next error")
	err = errors.WrapTag(err,
		errors.T("event", "test event"),
		errors.T("test123", 123),
		errors.T("test", "hello"),
		errors.T("fn_stack", stack.CallerWithFunc(stack.CallerWithFunc)),
	)

	err = errors.WrapStack(err)
	err = errors.Wrapf(err, "next error name=%s", "wrapf")
	errors.Debug(err)
}
