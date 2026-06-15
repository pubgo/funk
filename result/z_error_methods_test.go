package result_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
)

func TestErrorWithFn(t *testing.T) {
	assert.True(t, result.Error{}.WithFn(func() error { return nil }).IsOK())
	assert.True(t, result.Error{}.WithFn(func() error { return errors.New("fn") }).IsErr())
	assert.True(t, result.Error{}.WithFn(nil).IsErr())
}

func TestErrorWithErrorfAndMapErr(t *testing.T) {
	err := result.Error{}.WithErrorf("formatted %s", "err")
	assert.True(t, err.IsErr())
	assert.Contains(t, err.Err().Error(), "formatted")

	mapped := result.ErrOf(errors.New("root")).MapErr(func(err error) error {
		return fmt.Errorf("mapped: %w", err)
	})
	assert.Contains(t, mapped.Err().Error(), "mapped")
}

func TestErrorMatchWithError(t *testing.T) {
	err := result.ErrOf(errors.New("root")).MatchWithError(
		func() error { return nil },
		func(err error) error { return fmt.Errorf("wrapped: %w", err) },
	)
	assert.Contains(t, err.Error(), "wrapped")

	ok := result.ErrOf(nil).MatchWithError(
		func() error { return nil },
		func(err error) error { return err },
	)
	assert.NoError(t, ok)
}

func TestErrorMustAndUnwrap(t *testing.T) {
	assert.NotPanics(t, func() {
		result.ErrOf(nil).Must()
	})

	defer recovery.Testing(t)
	assert.Panics(t, func() {
		result.ErrOf(errors.New("must")).Must()
	})

	assert.NotPanics(t, func() {
		_ = result.ErrOf(nil).Unwrap()
	})
	assert.Panics(t, func() {
		_ = result.ErrOf(errors.New("unwrap")).Unwrap()
	})
}

func TestErrorExpectAndMustWithLog(t *testing.T) {
	defer recovery.Testing(t)

	assert.Panics(t, func() {
		result.ErrOf(errors.New("expect")).Expect("failed %s", "now")
	})

	assert.NotPanics(t, func() {
		result.ErrOf(nil).MustWithLog()
	})
	assert.Panics(t, func() {
		result.ErrOf(errors.New("log")).MustWithLog()
	})
}

func TestErrorGettersAndInspect(t *testing.T) {
	err := result.ErrOf(errors.New("getter"))
	assert.True(t, err.IsErr())
	assert.Equal(t, "getter", err.GetErr().Error())
	assert.Equal(t, "getter", err.Err().Error())
	assert.Contains(t, err.String(), "getter")

	called := false
	err.InspectErr(func(error) { called = true })
	assert.True(t, called)

	assert.Equal(t, "OK", result.ErrOf(nil).String())
}
