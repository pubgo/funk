package errors

import (
	"fmt"
	"os"
)

func Must(err error, args ...interface{}) {
	if err == nil {
		return
	}

	err = WrapStack(Wrap(err, fmt.Sprint(args...)))
	Debug(err)
	panic(err)
}

func MustF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	err = WrapStack(Wrap(err, fmt.Sprintf(msg, args...)))
	Debug(err)
	panic(err)
}

func Must1[T any](ret T, err error) T {
	if err != nil {
		err = WrapStack(err)
		Debug(err)
		panic(err)
	}

	return ret
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
