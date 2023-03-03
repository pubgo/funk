package errors

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pubgo/funk/errors/internal"
	"unsafe"

	"github.com/alecthomas/repr"
	"github.com/rs/zerolog"

	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

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

func stringify(buf *bytes.Buffer, err error) {
	if err == nil {
		return
	}

	if err1, ok := err.(fmt.Stringer); !ok {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrMsg, err.Error()))
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrDetail, fmt.Sprintf("%#v", err)))
		err = Unwrap(err)
		if err != nil {
			buf.WriteString("====================================================================\n")
			stringify(buf, err)
		}
	} else {
		buf.WriteString("====================================================================\n")
		buf.WriteString(err1.String())
	}
}

func jsonify(err error) map[string]any {
	if err == nil {
		return nil
	}

	return nil
}
