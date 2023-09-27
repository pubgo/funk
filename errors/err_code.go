package errors

import (
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func NewCodeErr(code *errorpb.ErrCode) error {
	if generic.IsNil(code) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Err:    parseToProto(code),
			Caller: stack.Caller(1).String(),
		},
	}
}

func WrapCode(err error, code *errorpb.ErrCode) error {
	if generic.IsNil(err) {
		return nil
	}

	if code == nil {
		panic("error code is nil")
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Err:    parseToProto(code),
			Caller: stack.Caller(1).String(),
		},
	}
}
