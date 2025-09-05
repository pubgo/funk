package errinter

import (
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Error interface {
	ID() string
	Kind() string
	Error() string
	String() string
	MarshalJSON() ([]byte, error)
}

type ErrUnwrap interface {
	Unwrap() error
}

type ErrorProto interface {
	Proto() proto.Message
}

type GRPCStatus interface {
	GRPCStatus() *status.Status
}
