package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/repr"

	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func handleGrpcError(err error) error {
	if err == nil {
		return nil
	}

	switch v := err.(type) {
	case *ErrWrap:
		return v

	case GRPCStatus:
		return NewCodeErr(&errorpb.ErrCode{
			Message:    v.GRPCStatus().Message(),
			StatusCode: errorpb.Code(v.GRPCStatus().Code()),
			Name:       "lava.grpc.status",
			Details:    v.GRPCStatus().Proto().Details,
		})
	default:
		return err
	}
}

func parseError(val interface{}) error {
	if generic.IsNil(val) {
		return nil
	}

	switch v := val.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	case []byte:
		return errors.New(convert.B2S(v))
	default:
		return &Err{Msg: fmt.Sprintf("%v", v), Detail: repr.String(v)}
	}
}

func errStringify(buf *bytes.Buffer, err error) {
	if err == nil {
		return
	}

	err1, ok := err.(fmt.Stringer)
	if ok {
		if _, ok = err.(*ErrWrap); !ok {
			buf.WriteString("error:\n")
		}
		buf.WriteString(err1.String())
		return
	}

	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrMsg, strings.TrimSpace(err.Error())))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrDetail, strings.TrimSpace(fmt.Sprintf("%v", err))))
	err = Unwrap(err)
	if err != nil {
		errStringify(buf, err)
	}
}

func errJsonify(err error) map[string]any {
	if err == nil {
		return make(map[string]any)
	}

	var data = make(map[string]any, 6)
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

func strFormat(f fmt.State, verb rune, err Error) {
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

		if _, ok := skipStackMap.Load(cc.Pkg); ok {
			continue
		}

		ss = append(ss, cc)
	}
	return ss
}
