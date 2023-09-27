package errors

import (
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func NewRedirectErr(err *errorpb.ErrRedirect) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		pb: &errorpb.ErrWrap{
			Err:    parseToProto(err),
			Caller: stack.Caller(1).String(),
		},
	}
}
