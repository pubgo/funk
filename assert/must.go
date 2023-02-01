package assert

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/try"
)

func Must(err error, args ...interface{}) {
	if errors.IsNil(err) {
		return
	}

	panic(errors.WrapCaller(errors.Wrap(err, fmt.Sprint(args...))))
}

func MustFn(errFn func() error, args ...interface{}) {
	var err = try.Try(errFn)
	if errors.IsNil(err) {
		return
	}

	panic(errors.WrapCaller(errors.Wrap(err, fmt.Sprint(args...))))
}

func MustF(err error, msg string, args ...interface{}) {
	if errors.IsNil(err) {
		return
	}

	panic(errors.WrapCaller(errors.Wrap(err, fmt.Sprintf(msg, args...))))
}

func Must1[T any](ret T, err error) T {
	if !errors.IsNil(err) {
		panic(errors.WrapCaller(err))
	}

	return ret
}

func Exit(err error, args ...interface{}) {
	if errors.IsNil(err) {
		return
	}

	errors.Debug(errors.WrapCaller(errors.Wrap(err, fmt.Sprint(args...))))
	debug.PrintStack()
	os.Exit(1)
}

func ExitFn(errFn func() error, args ...interface{}) {
	var err = try.Try(errFn)
	if errors.IsNil(err) {
		return
	}

	errors.Debug(errors.WrapCaller(errors.Wrap(err, fmt.Sprint(args...))))
	debug.PrintStack()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if errors.IsNil(err) {
		return
	}

	errors.Debug(errors.WrapCaller(errors.Wrapf(err, msg, args...)))
	debug.PrintStack()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if !errors.IsNil(err) {
		errors.Debug(errors.WrapCaller(err))
		debug.PrintStack()
		os.Exit(1)
	}

	return ret
}
