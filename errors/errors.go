package errors

import (
	"errors"
	"fmt"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target) //nolint
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Cause(err error) error {
	for {
		err1 := errors.Unwrap(err)
		if err1 == nil {
			return err
		}

		err = err1
	}
}

func Wrap(err error, args ...interface{}) error {
	if err == nil || isNil(err) {
		return nil
	}

	return parseXErr(err, func(err *errImpl) { err.msg = fmt.Sprint(args...) })
}

func WrapFn(err error, fn func(err XErr) XErr) error {
	if err == nil || isNil(err) {
		return nil
	}

	return fn(parseXErr(err))
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil || isNil(err) {
		return nil
	}

	return parseXErr(err, func(err *errImpl) { err.msg = fmt.Sprintf(format, args...) })
}

func Parse(err error) XErr {
	if err == nil || isNil(err) {
		return nil
	}

	return parseXErr(err)
}

func WrapMap(err error, m Map) error {
	if err == nil || isNil(err) {
		return nil
	}

	return parseXErr(err, func(err *errImpl) {
		if m == nil || len(m) == 0 {
			return
		}

		for k := range m {
			err.tags[k] = m[k]
		}
	})
}
