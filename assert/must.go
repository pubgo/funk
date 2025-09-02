package assert

import (
	"fmt"
	"log"
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

	log.Printf("[ERROR] %s", fmt.Sprint(args...))
	debug.PrintStack()
	os.Exit(1)
}

func ExitFn(errFn func() error, args ...interface{}) {
	err := try(errFn)
	if err == nil {
		return
	}

	log.Printf("[ERROR] %s", fmt.Sprint(args...))
	debug.PrintStack()
	os.Exit(1)
}

func ExitF(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	log.Printf("[ERROR] %s", fmt.Sprintf(msg, args...))
	debug.PrintStack()
	os.Exit(1)
}

func Exit1[T any](ret T, err error) T {
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		debug.PrintStack()
		os.Exit(1)
	}

	return ret
}
