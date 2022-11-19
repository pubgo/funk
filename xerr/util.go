package xerr

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/funk/stack"
)

func p(a ...interface{}) { _, _ = fmt.Fprintln(os.Stderr, a...) }

func WrapXErr(err error, fns ...func(err *XError)) XErr {
	if IsNil(err) {
		return nil
	}

	err1 := &XError{Err: err, Meta: make(map[string]interface{})}
	if e, ok := err.(*XError); ok {
		err1 = e
	} else {
		for i := 0; ; i++ {
			var cc = stack.CallerWithDepth(CallStackDepth + i)
			if cc == "" {
				break
			}
			err1.Caller = append(err1.Caller, cc)
		}
	}

	if len(fns) > 0 {
		fns[0](err1)
	}

	return err1
}

func trans(err error) *XError {
	if IsNil(err) {
		return nil
	}

	switch err := err.(type) {
	case *XError:
		return err
	case interface{ Unwrap() error }:
		if err.Unwrap() == nil {
			return &XError{Detail: fmt.Sprintf("%#v", err)}
		}
		return &XError{Err: err.Unwrap(), Msg: err.Unwrap().Error()}
	default:
		return &XError{Msg: err.Error(), Detail: fmt.Sprintf("%#v", err)}
	}
}

func IsNil(err error) bool {
	if err == nil {
		return true
	}

	var v = reflect.ValueOf(err)
	if !v.IsValid() {
		return true
	}

	return v.IsZero()
}
