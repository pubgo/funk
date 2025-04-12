package anyhow

import (
	"context"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/internal/anyhow/aherrcheck"
	"github.com/pubgo/funk/log"
)

func NewChecker(errSetter *error) *Checker {
	if errSetter == nil {
		debug.PrintStack()
		panic("errSetter is nil")
	}

	return &Checker{errSetter: errSetter}
}

type Checker struct {
	_ [0]func()

	errSetter *error
	ctx       context.Context
	args      any
}

func (c *Checker) SetCtx(ctx context.Context) *Checker {
	c.ctx = ctx
	return c
}

func (c *Checker) SetArgs(args any) *Checker {
	c.args = args
	return c
}

func (c *Checker) Check(err error) bool {
	if err == nil {
		return false
	}

	// err No checking, repeat setting
	if (*c.errSetter) != nil {
		log.Warn().Msgf("setter is not nil, err=%v", *c.errSetter)
	}

	var ctx = context.Background()
	if c.ctx != nil {
		ctx = c.ctx
	}

	for _, fn := range aherrcheck.GetErrChecks() {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	*c.errSetter = errors.WrapCaller(err, 1)
	return true
}
