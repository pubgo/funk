package errcheck

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/log"
	"github.com/samber/lo"
)

func RecoveryAndCheck(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		errMust(fmt.Errorf("setter is nil"))
		return
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
		errMust(fmt.Errorf("errSetter is nil"))
		return false
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

	return errors.WrapCaller(fn(err), 1)
}

func Inspect(err error, fn func(err error)) {
	if err == nil {
		return
	}

	fn(err)
}

func InspectLog(err error, fn func(logger *log.Event), contexts ...context.Context) {
	if err == nil {
		return
	}

	fn(log.Err(err, contexts...))
}

func LogErr(err error, contexts ...context.Context) {
	if err == nil {
		return
	}

	log.Err(err, contexts...).
		CallerSkipFrame(1).
		Msg(err.Error())
}
