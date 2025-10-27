package errors

import (
	"github.com/pubgo/funk/v2/internal/errors/errinter"
)

type (
	Maps         = errinter.Maps
	Tags         = errinter.Tags
	Tag          = errinter.Tag
	ErrIs        = errinter.ErrIs
	ErrAs        = errinter.ErrAs
	ErrUnwrapper = errinter.ErrUnwrapper
	Error        = errinter.Error
	ErrorProto   = errinter.ErrorProto
	GRPCStatus   = errinter.GRPCStatus
)
