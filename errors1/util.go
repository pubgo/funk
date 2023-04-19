package error1

import (
	"fmt"

	"github.com/alecthomas/repr"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func GetError(err *ErrWrap) *errorpb.Error {
	var pb = err.Proto()
	var errPb = new(errorpb.Error)
	for {
		if pb == nil {
			break
		}

		if pb.Err != nil {
			if errPb.Msg == nil && pb.GetMsg() != nil {
				errPb.Msg = pb.GetMsg()
			}

			if errPb.Code == nil && pb.GetCode() != nil {
				errPb.Code = pb.GetCode()
			}

			if errPb.Trace == nil && pb.GetTrace() != nil {
				errPb.Trace = pb.GetTrace()
			}
		}

		if errPb.Msg != nil && errPb.Code != nil && errPb.Trace != nil {
			break
		}

		pb = pb.Child
	}

	return errPb
}

func ParseErr(err error) *ErrWrap {
	switch e := err.(type) {
	case *ErrMsg:
		return &ErrWrap{
			err: e,
			wrap: &errorpb.ErrWrap{
				Child:  nil,
				Caller: stack.Caller(1).String(),
				Err:    &errorpb.ErrWrap_Msg{Msg: e.Proto()},
			},
		}
	case *ErrWrap:
		return &ErrWrap{
			err: e,
			wrap: &errorpb.ErrWrap{
				Child: e.Proto(),
			},
		}
	default:
		return &ErrWrap{
			err: err,
			wrap: &errorpb.ErrWrap{
				Child: nil,
				Err: &errorpb.ErrWrap_Msg{
					Msg: &errorpb.ErrMsg{
						Msg:    err.Error(),
						Detail: fmt.Sprintf("%v", err),
						Stack:  repr.String(err),
					},
				},
			},
		}
	}
}

func ParseVal(val interface{}) *ErrWrap {
	if generic.IsNil(val) {
		return nil
	}

	var pb = &errorpb.ErrMsg{
		Detail: fmt.Sprintf("%v", val),
		Stack:  repr.String(val),
	}

	switch e := val.(type) {
	case error:
		return ParseErr(e)
	case string:
		pb.Msg = e
	case []byte:
		pb.Msg = string(e)
	case fmt.Stringer:
		pb.Msg = e.String()
	default:
		pb.Msg = fmt.Sprint(e)
	}

	return &ErrWrap{
		err: &ErrMsg{msg: pb},
		wrap: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Err:    &errorpb.ErrWrap_Msg{Msg: pb},
		},
	}
}
