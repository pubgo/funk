// Code generated by protoc-gen-go-errors. DO NOT EDIT.
// versions:
// - protoc-gen-go-errors v0.0.4
// - protoc                 v4.24.3
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
	Name:       "demo.test.v1.ok",
	Reason:     "ok",
	StatusCode: errorpb.Code_OK,
}
var _ = errors.RegisterErrCodes(ErrCodeOK)

var ErrCodeNotFound = &errorpb.ErrCode{
	Code:       int32(100000),
	Name:       "demo.test.v1.not_found",
	Reason:     "not found 找不到",
	StatusCode: errorpb.Code_NotFound,
}
var _ = errors.RegisterErrCodes(ErrCodeNotFound)

var ErrCodeUnknown = &errorpb.ErrCode{
	Code:       int32(100001),
	Name:       "demo.test.v1.unknown",
	Reason:     "unknown 未知",
	StatusCode: errorpb.Code_NotFound,
}
var _ = errors.RegisterErrCodes(ErrCodeUnknown)

var ErrCodeDbConn = &errorpb.ErrCode{
	Code:       int32(100003),
	Name:       "demo.test.v1.db_conn",
	Reason:     "db connect error",
	StatusCode: errorpb.Code_Internal,
}
var _ = errors.RegisterErrCodes(ErrCodeDbConn)

var ErrCodeUnknownCode = &errorpb.ErrCode{
	Code:       int32(100004),
	Name:       "demo.test.v1.unknown_code",
	Reason:     "default code",
	StatusCode: errorpb.Code_Internal,
}
var _ = errors.RegisterErrCodes(ErrCodeUnknownCode)

var ErrCodeCustomCode = &errorpb.ErrCode{
	Code:       int32(100005),
	Name:       "demo.test.v1.custom_code",
	Reason:     "this is custom msg",
	StatusCode: errorpb.Code_OK,
}
var _ = errors.RegisterErrCodes(ErrCodeCustomCode)
