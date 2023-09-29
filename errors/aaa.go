package errors

import (
	"google.golang.org/grpc/status"
)

type GRPCStatus interface {
	GRPCStatus() *status.Status
}

type ErrEqual interface {
	IsEqual(any) bool
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type ErrUnwrap interface {
	Unwrap() error
}
