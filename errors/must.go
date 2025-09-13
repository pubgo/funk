package errors

import (
	"fmt"
	"os"

	"github.com/pubgo/funk/assert"
)

func Must(err error, args ...interface{}) {
	assert.Must(err, args...)
}

func MustF(err error, msg string, args ...interface{}) {
	assert.MustF(err, msg, args...)
}

func Must1[T any](ret T, err error) T {
	return assert.Must1(ret, err)
}

func Exit(err error, args ...interface{}) {
	if err == nil {
		return
	}

	Debug(WrapStack(Wrap(err, fmt.Sprint(args...))))
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	Debug(WrapStack(Wrapf(err, msg, args...)))
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if err != nil {
		Debug(WrapStack(err))
		os.Exit(1)
	}

	return ret
}
