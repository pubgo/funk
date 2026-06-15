package result_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
)

func TestResultMarshalJSON(t *testing.T) {
	okBytes, err := json.Marshal(result.OK(42))
	require.NoError(t, err)
	assert.Equal(t, "42", string(okBytes))

	fail := result.Fail[int](errors.New("boom"))
	_, err = json.Marshal(fail)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestErrorMarshalJSON(t *testing.T) {
	okBytes, err := json.Marshal(result.ErrOf(nil))
	require.NoError(t, err)
	assert.Equal(t, "null", string(okBytes))

	errBytes, err := json.Marshal(result.ErrOf(errors.New("marshal-me")))
	require.NoError(t, err)
	assert.Contains(t, string(errBytes), "marshal-me")
}

func TestFailNilPanics(t *testing.T) {
	defer recovery.Testing(t)

	assert.Panics(t, func() {
		_ = result.Fail[int](nil)
	})
}

func TestFlatMapAlias(t *testing.T) {
	r := result.OK("hello").
		FlatMap(func(s string) result.Result[string] {
			return result.OK(s + " world")
		})

	val, ok := r.TryUnwrap()
	assert.True(t, ok)
	assert.Equal(t, "hello world", val)

	short := result.OK("").
		FlatMap(func(s string) result.Result[string] {
			if s == "" {
				return result.Fail[string](fmt.Errorf("empty"))
			}
			return result.OK(s)
		})
	assert.True(t, short.IsErr())
}

func TestFlatMapToAlias(t *testing.T) {
	r := result.FlatMapTo(result.OK(2), func(v int) result.Result[string] {
		return result.OK(fmt.Sprintf("n=%d", v))
	})

	val, ok := r.TryUnwrap()
	assert.True(t, ok)
	assert.Equal(t, "n=2", val)
}

func TestErrorWithErrNilPreservesState(t *testing.T) {
	ok := result.ErrOf(nil).WithErr(nil)
	assert.True(t, ok.IsOK())

	err := result.ErrOf(errors.New("base")).WithErr(nil)
	assert.True(t, err.IsErr())
	assert.Equal(t, "base", err.Err().Error())
}

func TestWithErrorfDiscardsValue(t *testing.T) {
	r := result.OK(1).WithErrorf("forced")
	assert.True(t, r.IsErr())
	assert.Contains(t, r.Err().Error(), "forced")
}

func TestUnwrapErrSemantics(t *testing.T) {
	val, err := result.OK(42).UnwrapErr()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)

	_, err = result.Fail[int](errors.New("nope")).UnwrapErr()
	assert.Error(t, err)
}

func TestSetErrorClearsValueOnResult(t *testing.T) {
	var r result.Result[int] = result.OK(99)

	assert.True(t, result.ErrOf(errors.New("propagate")).Throw(&r))
	assert.True(t, r.IsErr())
	assert.Equal(t, 0, r.UnwrapOrEmpty())
}

func TestApplyErrIgnoresTypedNilSetter(t *testing.T) {
	var r *result.Result[int]

	assert.False(t, result.ErrOf(errors.New("noop")).Throw(r))
}
