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

func Expect(err error, format string, args ...any) {
	if err == nil {
		return
	}

	err = errors.WrapCaller(err, 1)
	err = errors.Wrapf(err, format, args...)
	errMust(err)
}

func Map(err error, fn func(err error) error) error {
	if err == nil {
		return nil
	}

	return fn(err)
}

func InspectLog(err error, fn func(evt *log.Event), contexts ...context.Context) {
	if err == nil {
		return
	}

	fn(log.Err(err, contexts...))
}

func Inspect(err error, fn func(err error)) {
	if err == nil {
		return
	}

	fn(err)
}

func RecordLog(err error, contexts ...context.Context) {
	if err == nil {
		return
	}

	log.Err(err, contexts...).
		CallerSkipFrame(1).
		Msg(err.Error())
}
