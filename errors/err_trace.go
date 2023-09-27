package errors

import (
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func NewTraceErr(trace *errorpb.ErrTrace) error {
	if generic.IsNil(trace) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Err:    parseToProto(trace),
			Caller: stack.Caller(1).String(),
		},
	}
}

func WrapTrace(err error, trace *errorpb.ErrTrace) error {
	if generic.IsNil(err) {
		return nil
	}

	if trace == nil {
		panic("error trace is nil")
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Wrap:   parseErrToWrap(err),
			Err:    parseToProto(trace),
			Caller: stack.Caller(1).String(),
		},
	}
}
