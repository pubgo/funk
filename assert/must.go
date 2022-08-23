package assert

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/funk/xerr"
)

func Must(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprint(args...) }))
}

func MustF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	panic(xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprintf(msg, args...) }))
}

func Must1[T any](ret T, err error) T {
	if isErrNil(err) {
		return ret
	}

	panic(xerr.WrapXErr(err))
}

func Exit(err error, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprint(args...) }).DebugPrint()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if isErrNil(err) {
		return
	}

	xerr.WrapXErr(err, func(err *xerr.XError) { err.Detail = fmt.Sprintf(msg, args...) }).DebugPrint()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if isErrNil(err) {
		return ret
	}

	xerr.WrapXErr(err).DebugPrint()
	os.Exit(1)
	return ret
}

func isErrNil(err error) bool {
	return err == nil || reflect.ValueOf(err).IsNil()
}
