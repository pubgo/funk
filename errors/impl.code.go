package errors

import (
	"github.com/pubgo/funk/proto/errorpb"
)

func (t *baseErr) Code() errorpb.Code {
	return t.code
}

func (t *baseErr) AddCode(code errorpb.Code) {
	t.code = code
}
