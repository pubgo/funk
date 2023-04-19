package error1

import (
	"errors"
	"reflect"

	"github.com/pubgo/funk/generic"
)

func IfErr(err error, fn func(err error)) {
	if generic.IsNil(err) {
		return
	}

	fn(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func As(err error, target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}

	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	targetType := typ.Elem()
	if targetType.Kind() != reflect.Interface && !targetType.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}

	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}

		if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
			return true
		}

		err = Unwrap(err)
	}
	return false
}

func Unwrap(err error) error {
	u, ok := err.(ErrUnwrap)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func Cause(err error) error {
	for {
		err1 := Unwrap(err)
		if generic.IsNil(err1) {
			return err
		}

		err = err1
	}
}
