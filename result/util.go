package result

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"strings"
	"unsafe"

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

	if r.IsOK() {
		return false
	}

	isErr := func() bool {
		if setter != nil {
			return setter.IsErr()
		}

		if rawSetter != nil {
			return (*rawSetter) != nil
		}

		return false
	}

	getErr := func() error {
		if setter != nil {
			return setter.GetErr()
		}

		if rawSetter != nil {
			return *rawSetter
		}

		return nil
	}

	setErr := func(err error) {
		if setter != nil {
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

// errRecovery handles error recovery from panics
// This function is used to recover from panics and convert them to errors.
// It applies callback functions to transform the error if needed.
//
// Parameters:
//
//	getErr - A function that returns the current error (if any)
//	callbacks - Optional functions to transform the error
//
// Returns:
//
//	error - The recovered error, or nil if no error occurred
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
	if setter1 == nil && setter2 == nil {
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
		if err == nil {
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

	if setter == nil {
		panicIfError(errors.Errorf("error setter is nil"))
		return
	}

	switch errSet := setter.(type) {
	case *Error:
		errSet.err = err
	case *ProxyErr:
		*errSet.err = err
	case ProxyErr:
		*errSet.err = err
	default:
		// Use reflection for generic Result[T] types
		rv := reflect.ValueOf(setter)
		if !rv.IsValid() || rv.IsNil() {
			slog.Error("error setter is invalid or nil",
				slog.String("type", fmt.Sprintf("%T", setter)),
				slog.String("stack", string(debug.Stack())),
			)
			return
		}

		t := rv.Type()
		typeStr := t.String()

		// Check if it's a Result type (pointer or value)
		if !strings.Contains(typeStr, "Result[") {
			slog.Error("error setter type error, type is not Result",
				slog.String("type", fmt.Sprintf("%T", setter)),
				slog.String("type-string", typeStr),
				slog.String("stack", string(debug.Stack())),
			)
			return
		}

		// Handle both *Result[T] and Result[T]
		var resultPtr *Result[any]
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				slog.Error("error setter is nil pointer",
					slog.String("type", typeStr),
					slog.String("stack", string(debug.Stack())),
				)
				return
			}
			resultPtr = (*Result[any])(rv.UnsafePointer())
		} else {
			// For value types, get address
			if !rv.CanAddr() {
				slog.Error("error setter cannot get address",
					slog.String("type", typeStr),
					slog.String("stack", string(debug.Stack())),
				)
				return
			}
			resultPtr = (*Result[any])(unsafe.Pointer(rv.UnsafeAddr()))
		}

		if resultPtr != nil {
			resultPtr.err = err
		}
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
