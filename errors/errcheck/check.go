package errcheck

import (
	"context"
	"fmt"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/log"
	"github.com/samber/lo"
)

func RecoveryAndCheck(setter *error, callbacks ...func(err error) error) {
	if setter == nil {
		assert.Must(fmt.Errorf("errcheck: setter is nil"))
		return
	}

	err := errinter.ParseError(recover())
	if err != nil {
		err = errors.WrapStack(err)
	}

	gErr := *setter
	if err == nil {
		err = gErr
	}

	if err == nil {
		return
	}

	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	*setter = err
}

func Check(errSetter *error, err error, contexts ...context.Context) bool {
	defer func() {
		if *errSetter == nil {
			return
		}

		logErr(lo.FirstOr(contexts, nil), *errSetter)
	}()

	if errSetter == nil {
		assert.Must(fmt.Errorf("errcheck: errSetter is nil"))
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
	for _, fn := range getCheckersFromCtx(ctx) {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	*errSetter = errors.WrapCaller(err, 1)
	return true
}

func Must(err error, args ...any) {
	logErr(nil, err)
	errNilOrPanic(err, args...)
}

func Must1[T any](v T, err error) T {
	logErr(nil, err)
	errNilOrPanic(err)
	return v
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

func InspectLog(err error, fn func(evt *log.Event), contexts ...context.Context) {
	if err == nil {
		return
	}

	logErr(lo.FirstOr(contexts, context.Background()), err, fn)
}

func LogErr(err error, contexts ...context.Context) {
	if err == nil {
		return
	}

	logErr(lo.FirstOr(contexts, context.Background()), err)
}
