package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func IfErr(err error, fn func(err error) error) error {
	if generic.IsNil(err) {
		return nil
	}

	return fn(err)
}

func New(msg string) error {
	return &Err{Msg: msg}
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

func WrapStack(err error) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		stack:  getStack(),
	}
}

func WrapCaller(err error, skip ...int) error {
	if generic.IsNil(err) {
		return nil
	}

	var depth = 1
	if len(skip) > 0 {
		depth += skip[0]
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(depth),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: Tags{T("msg", fmt.Sprintf(format, args...))},
	}
}

func Wrap(err error, msg string) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: Tags{T("msg", msg)},
	}
}

func T(k string, v any) Tag {
	return Tag{K: k, V: v}
}

func WrapTag(err error, tags ...Tag) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: tags,
	}
}

func WrapFn(err error, fn func() Tags) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: fn(),
	}
}

func WrapKV(err error, key string, value any) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: Tags{T(key, value)},
	}
}

func WrapMsg(err error, msg *errorpb.ErrMsg) error {
	if generic.IsNil(err) {
		return nil
	}

	if msg == nil {
		panic("error msg is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrMsg{pb: msg, err: err},
	}
}

func WrapCode(err error, code *errorpb.ErrCode) error {
	if generic.IsNil(err) {
		return nil
	}

	if code == nil {
		panic("error code is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrCode{pb: code, err: err},
	}
}

func WrapTrace(err error, trace *errorpb.ErrTrace) error {
	if generic.IsNil(err) {
		return nil
	}

	if trace == nil {
		panic("error trace is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrTrace{pb: trace, err: err},
	}
}

func Append(err error, errs ...error) error {
	if err == nil && len(errs) == 0 {
		return nil
	}

	if len(errs) == 0 {
		return &ErrWrap{
			err:    err,
			caller: stack.Caller(1),
		}
	}

	if err == nil {
		return &ErrWrap{
			err:    &errorsImpl{errs: errs},
			caller: stack.Caller(1),
		}
	}

	var errL []error
	switch err1 := err.(type) {
	case Errors:
		return err1.Append(errs...)
	default:
		errL = make([]error, 0, len(errs)+1)
		errL = append(errL, err1)
		errL = append(errL, errs...)
	}

	return &ErrWrap{
		err:    &errorsImpl{errs: errL},
		caller: stack.Caller(1),
	}
}
