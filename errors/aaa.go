package errors

import (
	"github.com/pubgo/funk/v2/errors/errinter"
	"golang.org/x/xerrors"
)

type Maps = errinter.Maps
type Tags = errinter.Tags
type Tag = errinter.Tag
type ErrIs = errinter.ErrIs
type ErrAs = errinter.ErrAs
type ErrUnwrap = errinter.ErrUnwrap
type Error = errinter.Error
type ErrorProto = errinter.ErrorProto
type GRPCStatus = errinter.GRPCStatus

var (
	Opaque = xerrors.Opaque
)
