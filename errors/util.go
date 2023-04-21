package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/repr"
	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func newErr(err error, skip ...int) *ErrMsg {
	var sk = 2
	if len(skip) > 0 {
		sk = sk + skip[0]
	}

	return &ErrMsg{
		ErrBase: &ErrBase{err: err},
		pb:      &errorpb.ErrMsg{},
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

func errStringify(buf *bytes.Buffer, err error) {
	if err == nil {
		return
	}

	if err1, ok := err.(fmt.Stringer); !ok {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrMsg, err.Error()))
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrDetail, fmt.Sprintf("%v", err)))
		err = Unwrap(err)
		if err != nil {
			errStringify(buf, err)
		}
	} else {
		buf.WriteString(err1.String())
	}
}

func errJsonify(err error) map[string]any {
	if err == nil {
		return nil
	}

	var data = make(map[string]any, 2)
	if _err, ok := err.(json.Marshaler); ok {
		data["cause"] = _err
	} else {
		data["err_msg"] = err.Error()
		data["err_detail"] = fmt.Sprintf("%v", err)
		err = Unwrap(err)
		if err != nil {
			data["cause"] = errJsonify(err)
		}
	}
	return data
}

func Format(f fmt.State, verb rune, err Error) {
	switch verb {
	case 'v':
		var data, err = err.MarshalJSON()
		if err != nil {
			fmt.Fprintln(f, err.Error())
		} else {
			fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		fmt.Fprintln(f, err.String())
	}
}

func getStack() []*stack.Frame {
	var ss []*stack.Frame
	for i := 0; ; i++ {
		var cc = stack.Caller(1 + i)
		if cc == nil {
			break
		}

		if cc.IsRuntime() {
			continue
		}

		if _, ok := skipStack.Load(cc.Pkg); ok {
			continue
		}

		ss = append(ss, cc)
	}
	return ss
}
