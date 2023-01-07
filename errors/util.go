package errors

import (
	"errors"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
	"reflect"
)

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

func parseXError(val interface{}) XError {
	switch _val := val.(type) {
	case nil:
		return nil
	case XError:
		return _val
	case error:
		return newErr(_val)
	case string:
		return newErr(errors.New(_val))
	case []byte:
		return newErr(errors.New(string(_val)))
	default:
		return newErr(errors.New(pretty.Sprint(_val)))
	}
}
