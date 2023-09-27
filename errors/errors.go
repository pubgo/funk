package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
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

func UnwrapEach(err error, call func(e error) bool) {
	if err == nil {
		return
	}

	for {
		if !call(err) {
			return
		}

		err1, ok := err.(ErrUnwrap)
		if !ok {
			return
		}

		err = err1.Unwrap()
	}
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
		err:    handleGrpcError(err),
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
		err:    handleGrpcError(err),
		caller: stack.Caller(depth),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    handleGrpcError(err),
		caller: stack.Caller(1),
		fields: Tags{T("msg", fmt.Sprintf(format, args...))},
	}
}

func Wrap(err error, msg string) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    handleGrpcError(err),
		caller: stack.Caller(1),
		fields: Tags{T("msg", msg)},
	}
}

func WrapMapTag(err error, tags Maps) error {
	if generic.IsNil(err) {
		return nil
	}

	if tags == nil {
		return err
	}

	return &ErrWrap{
		err:    handleGrpcError(err),
		caller: stack.Caller(1),
		fields: tags.Tags(),
	}
}

func WrapTag(err error, tags ...Tag) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    handleGrpcError(err),
		caller: stack.Caller(1),
		fields: tags,
	}
}

func WrapFn(err error, fn func() Tags) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    handleGrpcError(err),
		caller: stack.Caller(1),
		fields: fn(),
	}
}

func WrapKV(err error, key string, value any) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    handleGrpcError(err),
		caller: stack.Caller(1),
		fields: Tags{T(key, value)},
	}
}

func T(k string, v any) Tag {
	return Tag{K: k, V: v}
}
