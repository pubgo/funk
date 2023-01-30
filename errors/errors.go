package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func IfErr(err error, fn func(err error)) {
	if IsNil(err) {
		return
	}

	fn(err)
}

func New(format string, a ...interface{}) error {
	return &baseErr{
		err:    fmt.Errorf(format, a...),
		caller: stack.Caller(1),
	}
}

func Parse(val interface{}) XError {
	return parseXError(val)
}

func ParseResp(err error) *RespErr {
	if IsNil(err) {
		return nil
	}

	var rsp = &RespErr{Tags: make(map[string]any), Cause: err, Msg: err.Error()}
	for err != nil {
		switch _err := err.(type) {
		case XError:
			if rsp.Code == 0 {
				rsp.Code = _err.Code()
			}

			if rsp.BizCode == "" {
				rsp.BizCode = _err.BizCode()
			}

			if tags := _err.Tags(); tags != nil && len(tags) > 0 {
				for k, v := range tags {
					rsp.Tags[k] = v
				}
			}
		}

		err = Unwrap(err)
	}
	return rsp
}

func Debug(err error) {
	if IsNil(err) {
		return
	}

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
	u, ok := err.(errUnwrap)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func Cause(err error) error {
	for {
		err1 := Unwrap(err)
		if err1 == nil || IsNil(err1) {
			return err
		}

		err = err1
	}
}

func WrapStack(err error) error {
	if IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.AddStack()
	return base
}

func WrapFn(err error, fn func(xrr XError)) error {
	if IsNil(err) {
		return nil
	}

	if fn == nil {
		panic("[fn] should not be nil")
	}

	base := newErr(err)
	fn(base)
	return base
}

func WrapCaller(err error, skip ...int) error {
	if IsNil(err) {
		return nil
	}

	return newErr(err, skip...)
}

func Wrap(err error, msg string) error {
	if IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.AddMsg(msg)
	return base
}

func Wrapf(err error, format string, args ...interface{}) error {
	if IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.AddMsg(fmt.Sprintf(format, args...))
	return base
}

func WrapTags(err error, m Tags) error {
	if IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.AddTags(m)
	return base
}

func WrapCode(err error, code errorpb.Code) error {
	if IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.AddCode(code)
	return base
}

func WrapBizCode(err error, bizCode string) error {
	if IsNil(err) {
		return nil
	}

	base := newErr(err)
	base.AddBizCode(bizCode)
	return base
}

func Append(err error, errs ...error) Errors {
	switch err := err.(type) {
	case Errors:
		var errL = make([]error, 0, len(err)+len(errs))
		errL = append(errL, err...)
		errL = append(errL, errs...)
		return errL
	default:
		var errL = make([]error, 0, len(errs)+1)
		errL = append(errL, err)
		errL = append(errL, errs...)
		return errL
	}
}
