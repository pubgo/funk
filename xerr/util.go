package xerr

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/pubgo/funk/internal/utils"
)

func isErrNil(err error) bool { return err == nil || reflect.ValueOf(err).IsNil() }
func p(a ...interface{})      { _, _ = fmt.Fprintln(os.Stderr, a...) }

func ParseErr(err *error, val interface{}) {
	switch _val := val.(type) {
	case nil:
		return
	case error:
		*err = _val
	case string:
		*err = errors.New(_val)
	case []byte:
		*err = errors.New(string(_val))
	default:
		*err = fmt.Errorf("%#v", _val)
	}
	*err = WrapXErr(*err)
}

func WrapXErr(err error, fns ...func(err *XError)) XErr {
	if isErrNil(err) {
		return nil
	}

	err1 := &XError{Err: err, Meta: make(map[string]interface{})}
	if e, ok := err.(*XError); ok {
		err1 = e
	} else {
		for i := 0; ; i++ {
			var cc = utils.CallerWithDepth(CallStackDepth + i)
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
	if isErrNil(err) {
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
