package result

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errparser"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/log/logfields"
	"github.com/pubgo/funk/v2/result/resultchecker"
	"github.com/pubgo/funk/v2/stack"
)

var errFnIsNil = errors.New("[fn] is nil")

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errFnIsNil)
		return gErr
	}

	defer func() {
		if err := errparser.Parse(recover()); err != nil {
			gErr = errors.WrapStack(err)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
		}
	}()

	gErr = fn()
	return gErr
}

func tryResult[T any](fn func() Result[T]) (r Result[T]) {
	if fn == nil {
		return r.WithErr(errFnIsNil)
	}

	defer func() {
		var gErr error
		if err := errparser.Parse(recover()); err != nil {
			gErr = errors.WrapStack(err)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}

		r = r.WithErr(gErr)
	}()

	return fn()
}

func try1[T any](fn func() (T, error)) (t T, gErr error) {
	if fn == nil {
		return t, errors.WrapStack(errFnIsNil)
	}

	defer func() {
		if err := errparser.Parse(recover()); err != nil {
			gErr = errors.WrapStack(err)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}
	}()

	t, gErr = fn()
	return t, gErr
}

// panicIfError logs the error and panics
// This maintains backward compatibility with existing code that expects panics
func panicIfError(err error, events ...func(e Event)) {
	if err == nil {
		return
	}

	logErr(context.Background(), 1, err, events...)
	panic(err)
}

// catchErr handles error propagation from a Result Error to an error setter
// This function is responsible for propagating errors from a Result Error
// to various types of error setters (ErrSetter or raw error pointers).
// It applies error checkers and wraps the error before setting it.
//
// Parameters:
//
//	r - The Result Error containing the error to propagate
//	setter - An ErrSetter interface implementation
//	rawSetter - A raw error pointer
//	contexts - Optional context for error checking
//
// Returns:
//
//	bool - true if an error was set, false otherwise
func catchErr(r Error, setter ErrSetter, rawSetter *error, contexts ...context.Context) bool {
	if setter == nil && rawSetter == nil {
		panicIfError(errors.Errorf("error setter is nil"))
	}

	if isNilErrSetter(setter) && rawSetter == nil {
		return false
	}

	if r.IsOK() {
		return false
	}

	isErr := func() bool {
		if !isNilErrSetter(setter) {
			return setter.IsErr()
		}

		if rawSetter != nil {
			return (*rawSetter) != nil
		}

		return false
	}

	getErr := func() error {
		if !isNilErrSetter(setter) {
			return setter.GetErr()
		}

		if rawSetter != nil {
			return *rawSetter
		}

		return nil
	}

	setErr := func(err error) {
		if !isNilErrSetter(setter) {
			setError(setter, err)
		}

		if rawSetter != nil {
			setError(ErrProxyOf(rawSetter), err)
		}
	}

	ctx := lo.FirstOr(contexts, context.Background())
	if ctx == nil {
		ctx = context.Background()
	}

	// err No checking, repeat setting
	if isErr() {
		err := getErr()
		log.Err(err, ctx).Msgf("error setter has already set the error, err=%s", err.Error())
	}

	checkers := append(resultchecker.GetErrChecks(), resultchecker.GetCheckersFromCtx(ctx)...)
	err := r.getErr()
	for _, fn := range checkers {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	setErr(errors.WrapCaller(err, 2))

	return true
}

// errRecovery handles error recovery from panics when recover is called directly
// in the deferred function. Prefer result.Recovery for deferred panic handling.
func errRecovery(getErr func() error, callbacks ...func(err error) error) error {
	err := errparser.Parse(recover())
	if err == nil {
		err = getErr()
	}

	if err == nil {
		return nil
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return nil
		}
	}

	stack.Print()
	return err
}

// unwrapErr unwraps a Result and handles error propagation
// This function extracts the value from a Result while handling error propagation
// to error setters. It applies error checkers and returns the value or error.
//
// Parameters:
//
//	r - The Result to unwrap
//	setter1 - A raw error pointer
//	setter2 - An ErrSetter interface implementation
//	contexts - Optional context for error checking
//
// Returns:
//
//	T - The unwrapped value
//	error - Any error that occurred during unwrapping
func unwrapErr[T any](r Result[T], setter1 *error, setter2 ErrSetter, contexts ...context.Context) (T, error) {
	if setter1 == nil && isNilErrSetter(setter2) {
		panicIfError(fmt.Errorf("error setter is nil"))
	}

	ret := r.getValue()
	if r.IsOK() {
		return ret, nil
	}

	ctx := lo.FirstOr(contexts, context.Background())
	if ctx == nil {
		ctx = context.Background()
	}

	getPreErr := func() error {
		err := lo.FromPtr(setter1)
		if err == nil && !isNilErrSetter(setter2) {
			err = setter2.GetErr()
		}
		return err
	}
	if preErr := getPreErr(); preErr != nil {
		log.Err(preErr, ctx).Msgf("error setter has already set the error, err=%v", preErr)
	}

	err := r.getErr()
	checkers := append(resultchecker.GetErrChecks(), resultchecker.GetCheckersFromCtx(ctx)...)
	for _, fn := range checkers {
		err = fn(ctx, err)
		if err == nil {
			return ret, nil
		}
	}

	return ret, err
}

// setError sets an error on an ErrSetter
// This function handles setting an error on various types of error setters,
// including Error, ProxyErr, and generic Result types using reflection.
//
// Parameters:
//
//	setter - The ErrSetter to set the error on
//	err - The error to set
func setError(setter ErrSetter, err error) {
	if err == nil {
		return
	}

	if isNilErrSetter(setter) {
		panicIfError(errors.Errorf("error setter is nil"))
		return
	}

	setter.applyErr(err)
}

func isNilErrSetter(setter ErrSetter) bool {
	if setter == nil {
		return true
	}

	v := reflect.ValueOf(setter)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

var resultFile = stack.Caller(0)

// logErr logs an error with detailed context and stack trace
// This function provides comprehensive error logging with stack traces,
// error IDs, and other contextual information.
//
// Parameters:
//
//	ctx - The context for logging
//	skip - Number of stack frames to skip
//	err - The error to log
//	events - Optional functions to add additional log fields
func logErr(ctx context.Context, skip int, err error, events ...func(e Event)) {
	if err == nil {
		return
	}

	traces := lo.Filter(stack.Trace(), func(item *stack.Frame, index int) bool {
		return !item.IsRuntime() && item.Pkg != resultFile.Pkg
	})

	log.Error(ctx).
		Func(func(e *zerolog.Event) {
			e.Str(logfields.Module, "result")
			e.Strs(logfields.ErrorStack, lo.Map(traces, func(item *stack.Frame, index int) string { return item.String() }))
			e.Str(logfields.ErrorID, errors.GetErrorId(err))
			e.Str(logfields.ErrorDetail, fmt.Sprintf("%v", err))
			e.Str(zerolog.ErrorFieldName, err.Error())
			e.CallerSkipFrame(2 + skip)
		}).
		Func(func(e *zerolog.Event) {
			evt := Event{e}
			for _, fn := range events {
				fn(evt)
			}
		}).
		Msgf("%s\n%s", err.Error(), errors.JsonPrint(err))
}
