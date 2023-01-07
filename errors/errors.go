package errors

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kr/pretty"
	"github.com/pubgo/funk/stack"
	"google.golang.org/grpc/codes"
)

func Parse(err *Error, val interface{}) {
	switch _val := val.(type) {
	case nil:
		return
	case Error:
		*err = _val
	case error:
		*err = newErr(_val)
	case string:
		*err = newErr(errors.New(_val))
	case []byte:
		*err = newErr(errors.New(string(_val)))
	default:
		*err = newErr(errors.New(pretty.Sprint(_val)))
	}
}

func ParseResp(err error) *RespErr {
	if err == nil || isNil(err) {
		return nil
	}

	var rsp = &RespErr{Tags: make(map[string]any), Cause: err, Msg: err.Error()}
	for err != nil {
		switch _err := err.(type) {
		case ICodeWrap:
			if rsp.Code == 0 {
				rsp.Code = _err.Code()
			}
		case IBizCodeWrap:
			if rsp.BizCode == "" {
				rsp.BizCode = _err.BizCode()
			}
		case ITagWrap:
			for k, v := range _err.Tags() {
				rsp.Tags[k] = v
			}
		case IMsgWrap:
			rsp.Msg = _err.Msg()
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

	var base = &errStackImpl{baseErr: newErr(err)}
	for i := 0; ; i++ {
		var cc = stack.Caller(1 + i)
		if cc == nil {
			break
		}

		if cc.IsRuntime() {
			continue
		}

		base.stacks = append(base.stacks, cc)
	}
	return base
}

func Err(err error) error {
	if err == nil || isNil(err) {
		return nil
	}

	return newErr(err)
}

func Wrap(err error, msg string) error {
	if err == nil || isNil(err) {
		return nil
	}

	return &errMsgImpl{baseErr: newErr(err), msg: msg}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil || isNil(err) {
		return nil
	}

	return &errMsgImpl{baseErr: newErr(err), msg: fmt.Sprintf(format, args...)}
}

func WrapTags(err error, m Map) error {
	if err == nil || isNil(err) {
		return nil
	}

	return &errTagImpl{baseErr: newErr(err), tags: m}
}

func WrapCode(err error, code codes.Code) error {
	if err == nil || isNil(err) {
		return nil
	}

	return &errCodeImpl{baseErr: newErr(err), code: code}
}

func WrapBizCode(err error, bizCode string) error {
	if err == nil || isNil(err) {
		return nil
	}

	return &errBizCodeImpl{baseErr: newErr(err), bizCode: bizCode}
}

func isNil(err error) bool {
	if err == nil {
		return true
	}

	var v = reflect.ValueOf(err)
	if !v.IsValid() {
		return true
	}

	return v.IsZero()
}

func newErr(err error) *baseErr {
	return &baseErr{
		err:    err,
		caller: stack.Caller(2),
	}
}
