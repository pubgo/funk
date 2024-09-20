package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func MustProtoToAny(p proto.Message) *anypb.Any {
	if p == nil {
		return nil
	}

	pb, err := anypb.New(p)
	if err != nil {
		log.Err(err).Any("data", p).Msgf("failed to encode protobuf message to any message")
		return nil
	} else {
		return pb
	}
}

func ParseErrToPb(err error) proto.Message {
	if err == nil {
		return nil
	}

	switch err1 := err.(type) {
	case ErrorProto:
		return err1.Proto()
	case GRPCStatus:
		return err1.GRPCStatus().Proto()
	case proto.Message:
		return err1
	default:
		return &errorpb.ErrMsg{Msg: err.Error(), Detail: fmt.Sprintf("%v", err)}
	}
}

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

	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorErrMsg, strings.TrimSpace(err.Error())))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorErrDetail, strings.TrimSpace(fmt.Sprintf("%v", err))))
	err = Unwrap(err)
	if err != nil {
		errStringify(buf, err)
	}
}

func errJsonify(err error) map[string]any {
	if err == nil {
		return make(map[string]any)
	}

	data := make(map[string]any, 6)
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
		data, err := err.MarshalJSON()
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
		cc := stack.Caller(1 + i)
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
