package errors

import (
	"errors"
	"strings"
	"unsafe"

	"github.com/alecthomas/repr"
	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
	"github.com/rs/zerolog"
)

var pkgRoot string

func init() {
	tt := stack.Callers(1)
	if len(tt) == 0 {
		return
	}

	var files = strings.Split(tt[0].File, "/")
	pkgRoot = strings.Join(files[:len(files)-1], "/")
}

func convertEvent(evt *zerolog.Event) *event {
	return (*event)(unsafe.Pointer(evt))
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

func parseError(val interface{}) error {
	if generic.IsNil(val) {
		return nil
	}

	switch _val := val.(type) {
	case Error:
		return _val
	case error:
		return newErr(_val, 1)
	case string:
		return newErr(errors.New(_val), 1)
	case []byte:
		return newErr(errors.New(convert.B2S(_val)), 1)
	default:
		return newErr(errors.New(repr.String(_val)), 1)
	}
}
