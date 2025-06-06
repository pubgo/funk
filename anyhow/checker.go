package anyhow

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
)

func ErrChecker(setter *Error, contexts ...context.Context) *Checker {
	if setter == nil {
		errMust(errors.Errorf("error setter is nil"))
	}

	return &Checker{setter: setter, contexts: contexts}
}

func RawErrChecker(errSetter *error, contexts ...context.Context) *Checker {
	if errSetter == nil {
		errMust(errors.Errorf("raw error setter is nil"))
	}

	return &Checker{errSetter: errSetter, contexts: contexts}
}

type Checker struct {
	_ [0]func() // disallow ==

	errSetter *error
	setter    *Error
	contexts  []context.Context
	args      any
}

func (c *Checker) SetArgs(args any) *Checker {
	c.args = args
	return c
}

func (c *Checker) Recovery(callbacks ...func(err error) error) {
	if c.errSetter != nil {
		errRecovery(
			c.errSetter,
			func() bool { return c.errSetter != nil },
			func() error { return *c.errSetter },
			func(err error) error { return err },
			callbacks...,
		)
	} else if c.setter != nil {
		errRecovery(
			c.setter,
			func() bool { return c.setter.IsErr() },
			func() error { return c.setter.GetErr() },
			func(err error) Error { return newError(err) },
			callbacks...,
		)
	} else {
		err := errors.Parse(recover())
		if generic.IsNil(err) {
			return
		}

		for i := range callbacks {
			err = callbacks[i](err)
			if err == nil {
				return
			}
		}

		debug.PrintStack()
		errors.Debug(errors.WrapStack(err))
		os.Exit(1)
	}
}

func (c *Checker) Check(err Error) bool {
	return errTo(err, c.setter, c.errSetter, c.contexts...)
}

func (c *Checker) CheckErr(err error) bool {
	return errTo(newError(err), c.setter, c.errSetter, c.contexts...)
}
