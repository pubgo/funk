package anyhow

import (
	"context"
	"runtime/debug"

	"github.com/pubgo/funk/anyhow/aherrcheck"
	"github.com/pubgo/funk/errors"
	"github.com/rs/zerolog/log"
)

func newError(err error) Error {
	return Error{err: err}
}

type Error struct {
	_ [0]func() // disallow ==

	err error
}

func (r Error) OnErr(callbacks ...func(err error) error) Error {
	if r.IsErr() {
		var err = r.getErr()
		for _, fn := range callbacks {
			err = fn(err)
			if err == nil {
				return Error{}
			}
		}
		return Error{err: errors.WrapCaller(err, 1)}
	}

	return r
}

func (r Error) RawErrTo(setter *error, ctx ...context.Context) bool {
	return errTo(r, nil, setter, ctx...)
}

func (r Error) ErrTo(setter *Error, ctx ...context.Context) bool {
	return errTo(r, setter, nil, ctx...)
}

func (r Error) IsErr() bool {
	return r.getErr() != nil
}

func (r Error) IsErrNil() bool { return r.getErr() == nil }

func (r Error) GetErr() error {
	if !r.IsErr() {
		return nil
	}

	return errors.WrapCaller(r.getErr(), 1)
}

func (r Error) getErr() error {
	return r.err
}

func (r Error) Expect(format string, args ...any) {
	if !r.IsErr() {
		return
	}

	debug.PrintStack()
	err := errors.WrapCaller(r.getErr(), 1)
	err = errors.Wrapf(err, format, args...)
	errors.Debug(err)
	panic(err)
}

func errTo(r Error, setter *Error, rawSetter *error, contexts ...context.Context) bool {
	if setter == nil {
		debug.PrintStack()
		panic("setter is nil")
	}

	if !r.IsErr() {
		return false
	}

	var setterIsErr = func() bool {
		if setter != nil {
			return (*setter).IsErr()
		}

		if rawSetter != nil {
			return (*rawSetter) != nil
		}

		return false
	}

	var setterGetErr = func() error {
		if setter != nil {
			return (*setter).getErr()
		}

		if rawSetter != nil {
			return *rawSetter
		}

		return nil

	}

	// err No checking, repeat setting
	if setterIsErr() {
		log.Warn().Msgf("setter is not nil, err=%v", setterGetErr())
	}

	var ctx = context.Background()
	if len(contexts) > 0 {
		ctx = contexts[0]
	}

	var err = r.getErr()
	for _, fn := range aherrcheck.GetErrChecks() {
		err = fn(ctx, err)
		if err == nil {
			return false
		}
	}

	err = errors.WrapCaller(err, 2)
	if setter != nil {
		*setter = newError(err)
	} else if rawSetter != nil {
		*rawSetter = err
	}

	return true
}
