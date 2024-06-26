// Code generated by protoc-gen-go-errors. DO NOT EDIT.
// versions:
// - protoc-gen-go-errors v0.0.5
// - protoc                 v4.25.3
// source: testcodepb/test.proto

package testcodepb

import (
	errors "github.com/pubgo/funk/errors"
	errorpb "github.com/pubgo/funk/proto/errorpb"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

var ErrCodeOK = &errorpb.ErrCode{
	Code:       int32(0),
	Message:    "ok",
	Name:       "demo.test.v1.ok",
	StatusCode: errorpb.Code_OK,
}
var _ = errors.RegisterErrCodes(ErrCodeOK)

var ErrCodeNotFound = &errorpb.ErrCode{
	Code:       int32(100000),
	Message:    "not found 找不到",
	Name:       "demo.test.v1.not_found",
	StatusCode: errorpb.Code_NotFound,
}
var _ = errors.RegisterErrCodes(ErrCodeNotFound)

var ErrCodeUnknown = &errorpb.ErrCode{
	Code:       int32(100001),
	Message:    "unknown 未知",
	Name:       "demo.test.v1.unknown",
	StatusCode: errorpb.Code_NotFound,
}
var _ = errors.RegisterErrCodes(ErrCodeUnknown)

var ErrCodeDbConn = &errorpb.ErrCode{
	Code:       int32(100003),
	Message:    "db connect error",
	Name:       "demo.test.v1.db_conn",
	StatusCode: errorpb.Code_Internal,
}
var _ = errors.RegisterErrCodes(ErrCodeDbConn)

var ErrCodeUnknownCode = &errorpb.ErrCode{
	Code:       int32(100004),
	Message:    "default code",
	Name:       "demo.test.v1.unknown_code",
	StatusCode: errorpb.Code_Internal,
}
var _ = errors.RegisterErrCodes(ErrCodeUnknownCode)

var ErrCodeCustomCode = &errorpb.ErrCode{
	Code:       int32(100005),
	Message:    "this is custom msg",
	Name:       "demo.test.v1.custom_code",
	StatusCode: errorpb.Code_OK,
}
var _ = errors.RegisterErrCodes(ErrCodeCustomCode)
