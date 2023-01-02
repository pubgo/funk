package errors

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/funk/stack"
)

func p(a ...interface{}) { _, _ = fmt.Fprintln(os.Stderr, a...) }

func parseXErr(err error, fns ...func(err *errImpl)) XErr {
	if err == nil || isNil(err) {
		return nil
	}

	err1 := &errImpl{err: err, tags: make(map[string]interface{})}
	if e, ok := err.(*errImpl); ok {
		err1 = e
	}

	if err1.caller == nil {
		err1.caller = stack.Caller(2)
	}

	if len(err1.stackTrace) == 0 {
		for i := 0; ; i++ {
			var cc = stack.Caller(2 + i)
			if cc == nil {
				break
			}

			if cc.IsRuntime() {
				continue
			}

			err1.stackTrace = append(err1.stackTrace, cc)
		}
	}

	if len(fns) > 0 {
		fns[0](err1)
	}

	return err1
}

func trans(err error) *errImpl {
	if err == nil || isNil(err) {
		return nil
	}

	switch err := err.(type) {
	case *errImpl:
		return err
	case interface{ Unwrap() error }:
		if err.Unwrap() == nil {
			return &errImpl{msg: fmt.Sprintf("%#v", err)}
		}
		return &errImpl{err: err.Unwrap(), msg: err.Unwrap().Error()}
	default:
		return &errImpl{msg: err.Error(), tags: map[string]interface{}{"err_detail": fmt.Sprintf("%#v", err)}}
	}
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
