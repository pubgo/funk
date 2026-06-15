package result_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/result/resultchecker"
)

func TestThrowPanicsOnNilSetter(t *testing.T) {
	defer recovery.Testing(t)

	assert.Panics(t, func() {
		result.ErrOf(errors.New("boom")).Throw(nil)
	})
}

func TestThrowErrUsesRawPointer(t *testing.T) {
	var raw error
	assert.True(t, result.ErrOf(errors.New("raw")).ThrowErr(&raw))
	assert.Error(t, raw)
}

func TestThrowReturnsFalseForOKError(t *testing.T) {
	var r result.Result[int]
	assert.False(t, result.ErrOf(nil).Throw(&r))
	assert.True(t, r.IsOK())
}

func TestThrowPropagatesToErrorSetter(t *testing.T) {
	var e result.Error
	assert.True(t, result.ErrOf(errors.New("to error")).Throw(&e))
	assert.True(t, e.IsErr())
}

func TestThrowDoesNotOverwriteWhenCheckerSuppresses(t *testing.T) {
	t.Cleanup(func() {
		resultchecker.RemoveErrCheck(suppressErrCheck)
	})
	resultchecker.RegisterErrCheck(suppressErrCheck)

	var r result.Result[int]
	assert.False(t, result.ErrOf(errors.New("filtered")).Throw(&r))
	assert.True(t, r.IsOK())
}

func TestThrowUsesContextChecker(t *testing.T) {
	ctx := resultchecker.CreateCtx(context.Background(), []resultchecker.ErrChecker{suppressErrCheck})

	var r result.Result[int]
	assert.False(t, result.ErrOf(errors.New("filtered")).Throw(&r, ctx))
	assert.True(t, r.IsOK())
}

func TestThrowLogsWhenSetterAlreadyHasError(t *testing.T) {
	var r result.Result[int] = result.Fail[int](errors.New("existing"))

	assert.NotPanics(t, func() {
		assert.True(t, result.ErrOf(errors.New("next")).Throw(&r))
	})
	assert.True(t, r.IsErr())
}

func TestRecoveryCapturesPanic(t *testing.T) {
	var r result.Error
	func() {
		defer result.Recovery(&r)
		panic("boom")
	}()
	assert.True(t, r.IsErr())
}

func TestRecoveryErrCapturesPanic(t *testing.T) {
	var raw error
	func() {
		defer result.RecoveryErr(&raw)
		panic("boom")
	}()
	assert.Error(t, raw)
}

func TestRecoveryCallbackCanClearError(t *testing.T) {
	var r result.Error
	func() {
		defer result.Recovery(&r, func(err error) error { return nil })
		panic("boom")
	}()
	assert.True(t, r.IsOK())
}

func TestApplyErrOnTypedNilErrorSetter(t *testing.T) {
	var e *result.Error

	assert.False(t, result.ErrOf(errors.New("noop")).Throw(e))
}

func TestProxyErrApplyErrNilInnerPointer(t *testing.T) {
	var proxy result.ProxyErr
	assert.NotPanics(t, func() {
		result.ErrOf(errors.New("noop")).Throw(&proxy)
	})
}

func suppressErrCheck(context.Context, error) error { return nil }
