package errcheck

import (
	"context"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/samber/lo"
)

func Recovery(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	err := errors.Parse(recover())
	if err == nil {
		debug.PrintStack()
	}

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

	// err No checking, repeat setting
	if (*errSetter) != nil {
		log.Warn().Msgf("setter is not nil, err=%v", *errSetter)
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
		log.Warn().Msgf("setter is not nil, err=%v", *errSetter)
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
