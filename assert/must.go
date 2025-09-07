package assert

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
)

func Must(err error, args ...interface{}) {
	if err == nil {
		return
	}

	must(err, args...)
}

func MustFn(errFn func() error, args ...interface{}) {
	err := try(errFn)
	if err == nil {
		return
	}

	must(err, args...)
}

func MustF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	must(err, fmt.Sprintf(msg, args...))
}

func Must1[T any](ret T, err error, args ...any) T {
	if err != nil {
		must(err, args...)
	}

	return ret
}

func Exit(err error, args ...interface{}) {
	if err == nil {
		return
	}

	slog.Error("os exit with error", "err", err, "msg", fmt.Sprint(args...))
	debug.PrintStack()
	os.Exit(1)
}

func ExitFn(errFn func() error, args ...interface{}) {
	err := try(errFn)
	if err == nil {
		return
	}

	slog.Error("os exit with error func", "err", err, "msg", fmt.Sprint(args...))
	debug.PrintStack()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	slog.Error("os exit with error format", "err", err, "msg", fmt.Sprintf(msg, args...))
	debug.PrintStack()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if err != nil {
		slog.Error("os exit with error unwrap", "err", err)
		debug.PrintStack()
		os.Exit(1)
	}

	return ret
}
