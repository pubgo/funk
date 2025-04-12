package errcheck

import (
	"context"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/samber/lo"
)

func RecoveryAndCheck(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	err := errors.Parse(recover())
	gErr := *setter
	if err == nil && gErr == nil {
		return
	}

	if err == nil {
		err = gErr
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	*setter = errors.WrapCaller(err, 1)
}

func Check(errSetter *error, err error, contexts ...context.Context) bool {
	if errSetter == nil {
		debug.PrintStack()
		panic("errSetter is nil")
	}

	if err == nil {
		return false
	}

	if (*errSetter) != nil {
		log.Err(*errSetter).Msgf("errcheck: setter is not nil, err=%v", *errSetter)
		return true
	}

	var ctx = lo.FirstOr(contexts, context.Background())
	for _, fn := range GetCheckersFromCtx(ctx) {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	*errSetter = errors.WrapCaller(err, 1)
	return true
}

func CheckCtx(ctx context.Context, errSetter *error, err error, errCheckers ...ErrChecker) bool {
	if errSetter == nil {
		debug.PrintStack()
		panic("errSetter is nil")
	}

	if err == nil {
		return false
	}

	// err No checking, repeat setting
	if (*errSetter) != nil {
		log.Err(*errSetter, ctx).Msgf("setter is not nil, err=%v", *errSetter)
		return true
	}

	for _, fn := range append(GetCheckersFromCtx(ctx), errCheckers...) {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	*errSetter = errors.WrapCaller(err, 1)
	return true
}

func NewChecker(ctx context.Context, errSetter *error, errCheckers ...ErrChecker) *Checker {
	if errSetter == nil {
		debug.PrintStack()
		panic("errSetter is nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return &Checker{ctx: ctx, errSetter: errSetter, errCheckers: errCheckers}
}

type Checker struct {
	ctx         context.Context
	errSetter   *error
	errCheckers []ErrChecker
}

func (c *Checker) Check(err error, errCheckers ...ErrChecker) bool {
	if err == nil {
		return false
	}

	// err No checking, repeat setting
	if (*c.errSetter) != nil {
		log.Err(*c.errSetter, c.ctx).Msgf("setter is not nil, err=%v", *c.errSetter)
		return true
	}

	checkers := append(append(GetCheckersFromCtx(c.ctx), c.errCheckers...), errCheckers...)
	for _, fn := range checkers {
		err = fn(c.ctx, err)
		if err == nil {
			return false
		}
	}

	*c.errSetter = errors.WrapCaller(err, 1)
	return true
}

func (c *Checker) Recovery() {
	err := errors.Parse(recover())
	gErr := *c.errSetter
	if err == nil && gErr == nil {
		return
	}

	if err == nil {
		err = gErr
	}

	checkers := append(GetCheckersFromCtx(c.ctx), c.errCheckers...)
	for _, fn := range checkers {
		err = fn(c.ctx, err)
		if err == nil {
			return
		}
	}

	*c.errSetter = errors.WrapCaller(err, 1)
}
