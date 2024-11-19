// Code generated by protoc-gen-go-errors. DO NOT EDIT.
// versions:
// - protoc-gen-go-errors v0.0.7
// - protoc               v5.28.2
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

var TestErrCodeOK = &errorpb.ErrCode{
	Code:       int32(0),
	Message:    "ok",
	Name:       "demo.test.v1.ok",
	StatusCode: errorpb.Code_OK,
}
var _ = errors.RegisterErrCodes(TestErrCodeOK)

var TestErrCodeNotFound = &errorpb.ErrCode{
	Code:       int32(100000),
	Message:    "not found 找不到",
	Name:       "demo.test.v1.not_found",
	StatusCode: errorpb.Code_NotFound,
}
var _ = errors.RegisterErrCodes(TestErrCodeNotFound)

var TestErrCodeUnknown = &errorpb.ErrCode{
	Code:       int32(100001),
	Message:    "unknown 未知",
	Name:       "demo.test.v1.unknown",
	StatusCode: errorpb.Code_NotFound,
}
var _ = errors.RegisterErrCodes(TestErrCodeUnknown)

var TestErrCodeDbConn = &errorpb.ErrCode{
	Code:       int32(100003),
	Message:    "db connect error",
	Name:       "demo.test.v1.db_conn",
	StatusCode: errorpb.Code_Internal,
}
var _ = errors.RegisterErrCodes(TestErrCodeDbConn)

var TestErrCodeUnknownCode = &errorpb.ErrCode{
	Code:       int32(100004),
	Message:    "default code",
	Name:       "demo.test.v1.unknown_code",
	StatusCode: errorpb.Code_Internal,
}
var _ = errors.RegisterErrCodes(TestErrCodeUnknownCode)

var TestErrCodeCustomCode = &errorpb.ErrCode{
	Code:       int32(100005),
	Message:    "this is custom msg",
	Name:       "demo.custom.code",
	StatusCode: errorpb.Code_OK,
}
var _ = errors.RegisterErrCodes(TestErrCodeCustomCode)
