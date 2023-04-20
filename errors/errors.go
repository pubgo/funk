package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func NewEvent() *Event {
	return zerolog.Dict()
}

func IfErr(err error, fn func(err error)) {
	if generic.IsNil(err) {
		return
	}

	fn(err)
}

func New(format string, a ...interface{}) error {
	var err = fmt.Errorf(format, a...)
	return &baseErr{
		msg:    err.Error(),
		err:    err,
		caller: stack.Caller(1),
	}
}

func Parse(val interface{}) error {
	return parseError(val)
}

func Debug(err error) {
	if generic.IsNil(err) {
		return
	}

	err = parseError(err)
	if _err, ok := err.(fmt.Stringer); ok {
		fmt.Println(_err.String())
		return
	}

	pretty.Println(err)
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

func WrapStack(err error) error {
	if generic.IsNil(err) {
		return nil
	}

	var impl = &errStackImpl{err: err}
	impl.AddStack()
	return impl
}

func WrapCaller(err error, skip ...int) error {
	if generic.IsNil(err) {
		return nil
	}

	return newErr(err, skip...)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if generic.IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.msg = fmt.Sprintf(format, args...)
	return base
}

func Wrap(err error, msg string) error {
	if generic.IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.msg = msg
	return base
}

func WrapCode(err error, code ErrCode) error {
	if generic.IsNil(err) {
		return nil
	}

	if code == nil {
		panic("error code is nil")
	}

	return code.SetErr(err)
}

func WrapEventFn(err error, evt func(evt *Event)) error {
	if generic.IsNil(err) {
		return nil
	}

	base := &errEventImpl{err: err, caller: stack.Caller(1), evt: zerolog.Dict()}
	evt(base.evt)
	return base
}

func WrapEvent(err error, evt *Event) error {
	if generic.IsNil(err) {
		return nil
	}

	base := &errEventImpl{err: err, caller: stack.Caller(1), evt: evt}
	return base
}

func WrapKV(err error, k string, v any) error {
	if generic.IsNil(err) {
		return nil
	}

	var base ErrEvent
	switch err.(type) {
	case ErrEvent:
		base = err.(ErrEvent)
		base.Event().Any(k, v)
	default:
		base = &errEventImpl{err: err, caller: stack.Caller(1), evt: zerolog.Dict().Any(k, v)}
	}

	return base
}

func NewCode(code errorpb.Code) ErrCode {
	return &ErrCode{caller: stack.Caller(1), code: code, tags: make(map[string]string)}
}

func Append(err error, errs ...error) error {
	switch err1 := err.(type) {
	case Errors:
		var errL = make([]error, 0, len(err1.Errors())+len(errs))
		errL = append(errL, err1.Errors()...)
		errL = append(errL, errs...)
		return &errorsImpl{errs: errL}
	default:
		var errL = make([]error, 0, len(errs)+1)
		errL = append(errL, err1)
		errL = append(errL, errs...)
		return &errorsImpl{errs: errL}
	}
}
