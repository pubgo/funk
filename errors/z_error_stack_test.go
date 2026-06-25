package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/stack"
)

func TestStack(t *testing.T) {
	err := errors.WrapCaller(errors.New("hello error", errors.Tags{"name": "value"}))
	err = errors.Wrap(err, "next error")
	err = errors.WrapTags(err, errors.Tags{
		"event":    "test event",
		"test123":  123,
		"test":     "hello",
		"fn_stack": stack.CallerWithFunc(stack.CallerWithFunc),
	})

	err = errors.WrapStack(err)
	err = errors.Wrapf(err, "next error name=%s", "wrapf")

	wrapped, ok := errors.AsA[*errors.ErrWrap](err)
	assert.True(t, ok)

	var stacks []string
	errors.Walk(err, func(current error) bool {
		if wrap, ok := current.(*errors.ErrWrap); ok && len(wrap.Stacks) > 0 {
			stacks = wrap.Stacks
		}
		return true
	})
	assert.NotEmpty(t, stacks)
	assert.Contains(t, (*wrapped).String(), "hello error")
}
