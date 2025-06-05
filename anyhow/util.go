package anyhow

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/pubgo/funk/anyhow/aherrcheck"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/stack"
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

func must(err error, args ...interface{}) {
	if generic.IsNil(err) {
		return
	}

	if len(args) > 0 {
		err = errors.Wrap(err, fmt.Sprint(args...))
	}

	err = errors.WrapStack(err)
	errors.Debug(err)
	panic(err)
}

func errTo(r Error, setter *Error, rawSetter *error, contexts ...context.Context) bool {
	if setter == nil {
		must(errors.Errorf("error setter is nil"))
	}

	if r.IsOK() {
		return false
	}

	var isErr = func() bool {
		if setter != nil {
			return (*setter).IsErr()
		}

		if rawSetter != nil {
			return (*rawSetter) != nil
		}

		return false
	}

	var getErr = func() error {
		if setter != nil {
			return (*setter).getErr()
		}

		if rawSetter != nil {
			return *rawSetter
		}

		return nil
	}

	var setErr = func(err error) {
		if setter != nil {
			*setter = newError(err)
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

	var err = r.getErr()
	for _, fn := range aherrcheck.GetErrChecks() {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	setErr(errors.WrapCaller(err, 2))

	return true
}

func recovery[T any](setter *T, isErr func() bool, getErr func() error, newErr func(err error) T, callbacks ...func(err error) error) {
	if setter == nil {
		must(errors.Errorf("setter is nil"))
	}

	err := errors.Parse(recover())
	if err == nil && !isErr() {
		return
	}

	if err == nil {
		err = getErr()
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	err = errors.WrapCaller(err, 1)
	*setter = newErr(err)
}
