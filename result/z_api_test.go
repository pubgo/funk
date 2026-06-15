package result_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
)

func TestRunExecutesUntilFirstError(t *testing.T) {
	calls := 0
	err := result.Run(
		func() error {
			calls++
			return nil
		},
		func() error {
			calls++
			return errors.New("stop")
		},
		func() error {
			calls++
			return nil
		},
	)

	assert.True(t, err.IsErr())
	assert.Equal(t, 2, calls)
}

func TestAllCollectsValuesOrStopsOnError(t *testing.T) {
	ok := result.All(result.OK(1), result.OK(2), result.OK(3))
	vals, okUnwrap := ok.TryUnwrap()
	assert.True(t, okUnwrap)
	assert.Equal(t, []int{1, 2, 3}, vals)

	fail := result.All(result.OK(1), result.Fail[int](errors.New("broken")))
	assert.True(t, fail.IsErr())
}

func TestWrapAndWrapFn(t *testing.T) {
	ok := result.Wrap(7, nil)
	val, err := ok.UnwrapErr()
	assert.NoError(t, err)
	assert.Equal(t, 7, val)

	fail := result.Wrap(0, errors.New("wrap"))
	assert.True(t, fail.IsErr())

	fromFn := result.WrapFn(func() (int, error) { return 9, nil })
	assert.Equal(t, 9, fromFn.UnwrapOrEmpty())

	fromFnErr := result.WrapFn(func() (int, error) { return 0, errors.New("fn") })
	assert.True(t, fromFnErr.IsErr())
}

func TestWrapErr(t *testing.T) {
	val, err := result.WrapErr(5, nil)
	assert.True(t, err.IsOK())
	assert.Equal(t, 5, val)

	_, err = result.WrapErr(0, errors.New("wrap err"))
	assert.True(t, err.IsErr())
}

func TestErrOfFn(t *testing.T) {
	assert.True(t, result.ErrOfFn(func() error { return nil }).IsOK())
	assert.True(t, result.ErrOfFn(func() error { return errors.New("fn err") }).IsErr())
}

func TestMust1(t *testing.T) {
	assert.Equal(t, 11, result.Must1(11, nil))

	defer recovery.Testing(t)
	assert.Panics(t, func() {
		_ = result.Must1(0, errors.New("must1"))
	})
}

func TestPackageThrowHelpers(t *testing.T) {
	var r result.Result[string]
	assert.True(t, result.Throw(&r, errors.New("pkg throw")))
	assert.True(t, r.IsErr())

	var raw error
	assert.True(t, result.ThrowErr(&raw, errors.New("raw throw")))
	assert.Error(t, raw)
}

func TestErrProxyOfPanicsOnNilPointer(t *testing.T) {
	defer recovery.Testing(t)

	assert.Panics(t, func() {
		_ = result.ErrProxyOf(nil)
	})
}

func TestMapToPropagatesError(t *testing.T) {
	mapped := result.MapTo(result.Fail[int](errors.New("src")), func(v int) string {
		return fmt.Sprintf("%d", v)
	})
	assert.True(t, mapped.IsErr())

	ok := result.MapTo(result.OK(3), func(v int) string { return fmt.Sprintf("n=%d", v) })
	assert.Equal(t, "n=3", ok.UnwrapOrEmpty())
}

func TestCollectAndPartitionEdgeCases(t *testing.T) {
	empty := result.Collect([]result.Result[int]{})
	vals, ok := empty.TryUnwrap()
	assert.True(t, ok)
	assert.Empty(t, vals)

	values, errs := result.Partition([]result.Result[int]{
		result.OK(1),
		result.Fail[int](errors.New("e1")),
	})
	assert.Equal(t, []int{1}, values)
	assert.Len(t, errs, 1)
}

func TestLogErrHelpers(t *testing.T) {
	assert.NotPanics(t, func() {
		result.LogErr(errors.New("log me"))
		result.LogErrCtx(context.Background(), errors.New("log ctx"))
		result.LogErr(nil)
	})
}
