package assert

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/funk/errors"
)

func Must(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(errors.WrapXErr(err, func(err *errors.xerrImpl) { err.Detail = fmt.Sprint(args...) }))
}

func MustF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(errors.WrapXErr(err, func(err *errors.xerrImpl) { err.Detail = fmt.Sprintf(msg, args...) }))
}

func Must1[T any](ret T, err error) T {
	if isErrNil(err) {
		return ret
	}

	panic(errors.WrapXErr(err))
}

func Exit(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	errors.WrapXErr(err, func(err *errors.xerrImpl) { err.Detail = fmt.Sprint(args...) }).DebugPrint()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	errors.WrapXErr(err, func(err *errors.xerrImpl) { err.Detail = fmt.Sprintf(msg, args...) }).DebugPrint()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if isErrNil(err) {
		return ret
	}

	errors.WrapXErr(err).DebugPrint()
	os.Exit(1)
	return ret
}

func isErrNil(err error) bool {
	return err == nil || reflect.ValueOf(err).IsNil()
}
