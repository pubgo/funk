package errors

import (
	"errors"
	"reflect"

	"github.com/alecthomas/repr"
	"github.com/pubgo/funk/stack"
)

func IsNil(err interface{}) bool {
	if err == nil {
		return true
	}

	var v = reflect.ValueOf(err)
	if !v.IsValid() {
		return true
	}

	return v.IsZero()
}

func newErr(err error, skip ...int) *baseErr {
	var sk = 2
	if len(skip) > 0 {
		sk = sk + skip[0]
	}

	return &baseErr{
		err:    err,
		caller: stack.Caller(sk),
	}
}

func parseXError(val interface{}) XError {
	if IsNil(val) {
		return nil
	}

	switch _val := val.(type) {
	case XError:
		return _val
	case error:
		return newErr(_val, 1)
	case string:
		return newErr(errors.New(_val), 1)
	case []byte:
		return newErr(errors.New(string(_val)), 1)
	default:
		return newErr(errors.New(repr.String(_val)), 1)
	}
}
