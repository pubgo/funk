package errors

import (
	"errors"
	"fmt"

	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

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
	if err == nil || isNil(err) {
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
	fmt.Println(err.(fmt.Stringer).String())
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target) //nolint
}

//func Opaque(err error) error

func Unwrap(err error) error {
	u, ok := err.(errUnwrap)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

//func Opaque(err error) error {
//	return &opaqueWrapper{err: err}
//}

func Cause(err error) error {
	for {
		err1 := Unwrap(err)
		if err1 == nil || isNil(err1) {
			return err
		}

		err = err1
	}
}

func WrapStack(err error) error {
	if err == nil || isNil(err) {
		return nil
	}

	base := newErr(Parse(err))
	base.AddStack()
	return base
}

func WrapFn(err error, fn func(xrr XError)) error {
	if err == nil || isNil(err) {
		return nil
	}

	if fn == nil {
		panic("[fn] should not be nil")
	}

	base := newErr(Parse(err))
	fn(base)
	return base
}

func WrapCaller(err error) error {
	if err == nil || isNil(err) {
		return nil
	}

	return newErr(err)
}

func Wrap(err error, msg string) error {
	if err == nil || isNil(err) {
		return nil
	}

	base := newErr(Parse(err))
	base.AddMsg(msg)
	return base
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil || isNil(err) {
		return nil
	}

	base := newErr(Parse(err))
	base.AddMsg(fmt.Sprintf(format, args...))
	return base
}

func WrapTags(err error, m Tags) error {
	if err == nil || isNil(err) {
		return nil
	}

	base := newErr(Parse(err))
	base.AddTags(m)
	return base
}

func WrapCode(err error, code codes.Code) error {
	if err == nil || isNil(err) {
		return nil
	}

	base := newErr(Parse(err))
	base.AddCode(code)
	return base
}

func WrapBizCode(err error, bizCode string) error {
	if err == nil || isNil(err) {
		return nil
	}

	base := newErr(Parse(err))
	base.AddBizCode(bizCode)
	return base
}
