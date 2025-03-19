package errcheck

import (
	"context"
	"runtime/debug"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
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

	// err No checking, repeat setting
	if (*errSetter) != nil {
		log.Error().Stack().Msgf("errcheck: setter is not nil, err=%v", *errSetter)
		return true
	}

	var ctx = context.Background()
	if len(contexts) > 0 {
		ctx = contexts[0]
	}

	for _, fn := range GetErrChecks() {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	*errSetter = errors.WrapCaller(err, 1)
	return true
}
