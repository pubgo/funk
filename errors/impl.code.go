package errors

import (
	"google.golang.org/grpc/codes"
)

func (t *baseErr) Code() codes.Code {
	return t.code
}

func (t *baseErr) AddCode(code codes.Code) {
	t.code = code
}
