package result

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/internal/anyhow/aherrcheck"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/stack"
	"github.com/samber/lo"
)

var errFnIsNil = errors.New("[fn] is nil")

func try(fn func() error) (gErr error) {
	if fn == nil {
		gErr = errors.WrapStack(errFnIsNil)
		return
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			debug.PrintStack()
			errors.Debug(gErr)
		}

		gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn).String())
	}()

	gErr = fn()
	return
}

func try1[T any](fn func() (T, error)) (t T, gErr error) {
	if fn == nil {
		return t, errors.WrapStack(errFnIsNil)
	}

	defer func() {
		if err := errors.Parse(recover()); !generic.IsNil(err) {
			gErr = errors.WrapStack(err)
			debug.PrintStack()
			errors.Debug(gErr)
		}

		if gErr != nil {
			gErr = errors.WrapKV(gErr, "fn_stack", stack.CallerWithFunc(fn))
		}
	}()

	t, gErr = fn()
	return
}

func errMust(err error, args ...any) {
	if err == nil {
		return
	}

	if len(args) > 0 {
		err = errors.Wrap(err, fmt.Sprint(args...))
	}

	err = errors.WrapStack(err)
	errors.Debug(err)
	panic(err)
}

func catchErr(r Error, setter ErrSetter, rawSetter *error, contexts ...context.Context) bool {
	if setter == nil && rawSetter == nil {
		errMust(errors.Errorf("error setter is nil"))
	}

	if r.IsOK() {
		return false
	}

	var isErr = func() bool {
		if setter != nil {
			return setter.IsErr()
		}

		if rawSetter != nil {
			return (*rawSetter) != nil
		}

		return false
	}

	var getErr = func() error {
		if setter != nil {
			return setter.GetErr()
		}

		if rawSetter != nil {
			return *rawSetter
		}

		return nil
	}

	var setErr = func(err error) {
		if setter != nil {
			setter.setError(err)
		}

		if rawSetter != nil {
			*rawSetter = err
		}
	}

	// err No checking, repeat setting
	if isErr() {
		err := getErr()
		log.Err(err).Msgf("error setter is has value, err=%s", err.Error())
	}

	var ctx = context.Background()
	for i := range contexts {
		if contexts[i] == nil {
			continue
		}
		ctx = contexts[i]
		break
	}

	checkers := append(aherrcheck.GetErrChecks(), aherrcheck.GetCheckersFromCtx(ctx)...)
	var err = r.getErr()
	for _, fn := range checkers {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	setErr(errors.WrapCaller(err, 2))

	return true
}

func errRecovery(isErr func() bool, getErr func() error, callbacks ...func(err error) error) error {
	err := errors.Parse(recover())
	if err == nil && !isErr() {
		return nil
	}

	if err == nil {
		err = getErr()
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return nil
		}
	}
	return err
}

func unwrapErr[T any](r Result[T], setter1 *error, setter2 ErrSetter, contexts ...context.Context) (T, error) {
	if setter1 == nil && setter2 == nil {
		debug.PrintStack()
		panic("Unwrap: error setter is nil")
	}

	var ret = r.getValue()
	if r.IsOK() {
		return ret, nil
	}

	var ctx = context.Background()
	if len(contexts) > 0 {
		ctx = contexts[0]
	}

	getSetterErr := func() error {
		err := lo.FromPtr(setter1)
		if err == nil {
			err = setter2.GetErr()
		}
		return err
	}
	setterErr := getSetterErr()
	if setterErr != nil {
		log.Error(ctx).Msgf("Unwrap: error setter has value, err=%v", setterErr)
	}

	var err = r.getErr()
	for _, fn := range aherrcheck.GetErrChecks() {
		err = fn(ctx, err)
		if err == nil {
			return ret, nil
		}
	}

	return ret, err
}
