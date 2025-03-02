package anyhow

import (
	"context"
	"runtime/debug"

	"github.com/pubgo/funk/anyhow/aherrcheck"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
)

func Recovery(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	gErr := *setter
	err := errors.Parse(recover())
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

func ErrTo(errSetter *error, err error, contexts ...context.Context) bool {
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
	}

	var ctx = context.Background()
	if len(contexts) > 0 {
		ctx = contexts[0]
	}

	for _, fn := range aherrcheck.GetErrChecks() {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	*errSetter = errors.WrapCaller(err, 1)
	return true
}
