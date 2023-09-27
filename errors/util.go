package errors

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/alecthomas/repr"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func parseErrToWrap(err error) *errorpb.ErrWrap {
	if err != nil {
		return nil
	}

	switch v := err.(type) {
	case *ErrWrap:
		return v.pb
	default:
		return &errorpb.ErrWrap{
			Err:    parseToProto(err),
			Caller: stack.Stack(2).String(),
		}
	}
}

func parseToProto(err interface{}) *anypb.Any {
	if generic.IsNil(err) {
		return nil
	}

	switch v := err.(type) {
	case proto.Message:
		return lo.Must1(anypb.New(v))
	case GRPCStatus:
		return lo.Must1(anypb.New(&errorpb.ErrCode{
			Reason:  v.GRPCStatus().Message(),
			Code:    errorpb.Code(v.GRPCStatus().Code()),
			Name:    "lava.grpc.status",
			Details: v.GRPCStatus().Proto().Details,
		}))
	case error:
		return lo.Must1(anypb.New(&errorpb.ErrMsg{
			Msg:    v.Error(),
			Detail: fmt.Sprintf("%v", v),
		}))
	case string:
		return lo.Must1(anypb.New(&errorpb.ErrMsg{
			Msg: v,
		}))
	default:
		return lo.Must1(anypb.New(&errorpb.ErrMsg{
			Msg:    fmt.Sprintf("%v", err),
			Detail: repr.String(err),
		}))
	}
}

func getStack() []string {
	var ss []string
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

		ss = append(ss, cc.String())
	}
	return ss
}
