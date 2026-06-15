package result_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
)

func TestResultMapAndValidate(t *testing.T) {
	doubled := result.OK(2).Map(func(v int) int { return v * 2 })
	assert.Equal(t, 4, doubled.UnwrapOrEmpty())

	validated := result.OK("ok").Validate(func(s string) error {
		if s == "" {
			return errors.New("empty")
		}
		return nil
	})
	assert.True(t, validated.IsOK())

	invalid := result.OK("").Validate(func(s string) error {
		if s == "" {
			return errors.New("empty")
		}
		return nil
	})
	assert.True(t, invalid.IsErr())
}

func TestResultMapErrVariants(t *testing.T) {
	wrapped := result.Fail[int](errors.New("root")).MapErr(func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})
	assert.Contains(t, wrapped.Err().Error(), "wrapped")

	replaced := result.Fail[int](errors.New("root")).MapErrOr(func(err error) result.Result[int] {
		return result.OK(99)
	})
	assert.Equal(t, 99, replaced.UnwrapOrEmpty())

	unchanged := result.OK(1).MapErr(func(err error) error { return err })
	assert.Equal(t, 1, unchanged.UnwrapOrEmpty())
}

func TestResultFallbackHelpers(t *testing.T) {
	assert.Equal(t, 5, result.Fail[int](errors.New("x")).Or(5).UnwrapOrEmpty())
	assert.Equal(t, 6, result.Fail[int](errors.New("x")).UnwrapOr(6))
	assert.Equal(t, 7, result.Fail[int](errors.New("x")).OrElse(func() int { return 7 }).UnwrapOrEmpty())
	assert.Equal(t, 8, result.Fail[int](errors.New("x")).UnwrapOrElse(func() int { return 8 }))
}

func TestResultMatchWithResult(t *testing.T) {
	doubled := result.OK(2).MatchWithResult(
		func(v int) result.Result[int] { return result.OK(v * 2) },
		func(err error) result.Result[int] { return result.Fail[int](err) },
	)
	assert.Equal(t, 4, doubled.UnwrapOrEmpty())

	fromErr := result.Fail[int](errors.New("bad")).MatchWithResult(
		func(v int) result.Result[int] { return result.OK(v) },
		func(err error) result.Result[int] { return result.OK(-1) },
	)
	assert.Equal(t, -1, fromErr.UnwrapOrEmpty())
}

func TestResultWithFnAndWithValue(t *testing.T) {
	chained := result.OK(1).WithFn(func() (int, error) { return 2, nil })
	assert.Equal(t, 2, chained.UnwrapOrEmpty())

	skipped := result.Fail[int](errors.New("base")).WithFn(func() (int, error) {
		t.Fatal("should not run")
		return 0, nil
	})
	assert.True(t, skipped.IsErr())

	replaced := result.OK(1).WithValue(3)
	assert.Equal(t, 3, replaced.UnwrapOrEmpty())
}

func TestResultValueTo(t *testing.T) {
	var out int
	err := result.OK(10).ValueTo(&out)
	assert.True(t, err.IsOK())
	assert.Equal(t, 10, out)

	err = result.Fail[int](errors.New("bad")).ValueTo(&out)
	assert.True(t, err.IsErr())

	err = result.OK(1).ValueTo(nil)
	assert.True(t, err.IsErr())
}

func TestResultInspectHelpers(t *testing.T) {
	called := false
	result.OK(1).Inspect(func(int) { called = true })
	assert.True(t, called)

	seen := 0
	result.OK(1).IfOK(func(v int) { seen = v }).IfErr(func(error) { t.Fatal("unexpected") })
	assert.Equal(t, 1, seen)

	result.Fail[int](errors.New("x")).InspectErr(func(error) { called = true })
}

func TestResultCallIfOK(t *testing.T) {
	ok := result.OK(1).CallIfOK(func(int) error { return nil })
	assert.True(t, ok.IsOK())

	fail := result.OK(1).CallIfOK(func(int) error { return errors.New("call") })
	assert.True(t, fail.IsErr())

	unchanged := result.Fail[int](errors.New("base")).CallIfOK(func(int) error {
		t.Fatal("should not run")
		return nil
	})
	assert.True(t, unchanged.IsErr())
}

func TestResultStringAndWithErr(t *testing.T) {
	assert.Equal(t, "OK(7)", result.OK(7).String())
	assert.Contains(t, result.Fail[int](errors.New("bad")).String(), "Error")

	withTag := result.OK(1).WithErr(errors.New("tagged"), errors.Tags{"k": "v"})
	assert.True(t, withTag.IsErr())

	kept := result.OK(1).WithErr(nil)
	assert.True(t, kept.IsOK())
}

func TestResultExpectAndMust(t *testing.T) {
	assert.Equal(t, 3, result.OK(3).Expect("boom"))

	defer recovery.Testing(t)
	assert.Panics(t, func() {
		_ = result.Fail[int](errors.New("x")).Expect("custom %s", "msg")
	})
	assert.Panics(t, func() {
		result.Fail[int](errors.New("x")).Must()
	})
}

func TestResultUnwrapOrThrow(t *testing.T) {
	var r result.Result[int]
	val := result.OK(12).UnwrapOrThrow(&r)
	assert.Equal(t, 12, val)
	assert.True(t, r.IsOK())

	val = result.Fail[int](errors.New("throw")).UnwrapOrThrow(&r)
	assert.Equal(t, 0, val)
	assert.True(t, r.IsErr())
}
