package errors_test

import (
	"context"
	"testing"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/stack"
)

func TestErrF(t *testing.T) {
	err := errors.WrapCaller(errors.ErrF(func(ctx context.Context) error {

	}))
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
