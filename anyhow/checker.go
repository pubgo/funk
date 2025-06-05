package anyhow

import (
	"context"

	"github.com/pubgo/funk/errors"
	"github.com/samber/lo"
)

func ErrChecker(setter *Error, contexts ...context.Context) *Checker {
	if setter == nil {
		must(errors.Errorf("errSetter is nil"))
	}

	return &Checker{setter: setter, ctx: lo.FirstOr(contexts, nil)}
}

func RawErrChecker(errSetter *error, contexts ...context.Context) *Checker {
	if errSetter == nil {
		must(errors.Errorf("errSetter is nil"))
	}

	return &Checker{errSetter: errSetter, ctx: lo.FirstOr(contexts, nil)}
}

type Checker struct {
	_ [0]func() // disallow ==

	errSetter *error
	setter    *Error
	ctx       context.Context
	args      any
}

func (c *Checker) SetArgs(args any) *Checker {
	c.args = args
	return c
}

func (c *Checker) Recovery(callbacks ...func(err error) error) {

	if c.errSetter != nil {
		recovery(
			c.errSetter,
			func() bool { return c.errSetter != nil },
			func() error { return *c.errSetter },
			func(err error) error { return err },
			callbacks...,
		)
	}

	if c.setter != nil {
		recovery(
			c.setter,
			func() bool { return c.setter.IsErr() },
			func() error { return c.setter.GetErr() },
			func(err error) Error { return newError(err) },
			callbacks...,
		)
	}
}

func (c *Checker) Check(err error) bool {
	return errTo(newError(err), c.setter, c.errSetter, c.ctx)
}
