package anyhow

import (
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
func (r Error) ErrTo(setter *Error) bool {
	if setter == nil {
		debug.PrintStack()
		panic("ErrTo: setter is nil")
	}

	if !r.IsErr() {
		return false
	}

	// err No checking, repeat setting
	if (*setter).IsErr() {
		log.Warn().Msgf("ErrTo: setter is not nil, err=%v", (*setter).getErr())
	}

	var err = r.getErr()
	for _, fn := range aherrcheck.GetErrChecks() {
		err = fn(err)
		if err == nil {
			return false
		}
	}

	*setter = newError(errors.WrapCaller(err, 1))
	return true
}

func (r Error) Unwrap(setter *Error, callbacks ...func(err error) error) {
	if setter == nil {
		debug.PrintStack()
		panic("Unwrap: setter is nil")
	}

	if !r.IsErr() {
		return
	}

	// err No checking, repeat setting
	if (*setter).IsErr() {
		log.Warn().Msgf("Unwrap: setter is not nil, err=%v", (*setter).getErr())
	}

	callbacks = append(callbacks, aherrcheck.GetErrChecks()...)
	var err = r.getErr()
	for _, fn := range callbacks {
		err = fn(err)
		if err == nil {
			return
		}
	}

	err = errors.WrapCaller(err, 1)
	*setter = newError(err)
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
