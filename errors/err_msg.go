package errors

import (
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func NewMsgErr(msg *errorpb.ErrMsg) error {
	if generic.IsNil(msg) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Err:    parseToProto(msg),
			Caller: stack.Caller(1).String(),
		},
	}
}

func WrapMsg(err error, msg *errorpb.ErrMsg) error {
	if generic.IsNil(err) {
		return nil
	}

	if msg == nil {
		panic("error msg is nil")
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Err:    parseToProto(msg),
			Caller: stack.Caller(1).String(),
		},
	}
}
