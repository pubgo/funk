package errutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pubgo/funk/convert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/proto/errorpb"
)

func IsMemoryErr(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "invalid memory address or nil pointer dereference")
}

// Err2GrpcCode
// converts a standard Go error into its canonical code. Note that
// this is only used to translate the error returned by the server applications.
func Err2GrpcCode(err error) codes.Code {
	switch err {
	case nil:
		return codes.OK
	case io.EOF:
		return codes.OutOfRange
	case io.ErrClosedPipe, io.ErrNoProgress, io.ErrShortBuffer, io.ErrShortWrite, io.ErrUnexpectedEOF:
		return codes.FailedPrecondition
	case os.ErrInvalid:
		return codes.InvalidArgument
	case context.Canceled:
		return codes.Canceled
	case context.DeadlineExceeded:
		return codes.DeadlineExceeded
	}

	switch {
	case os.IsExist(err):
		return codes.AlreadyExists
	case os.IsNotExist(err):
		return codes.NotFound
	case os.IsPermission(err):
		return codes.PermissionDenied
	}
	return codes.Unknown
}

func Http2GrpcCode(code int32) codes.Code {
	switch code {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusPreconditionFailed:
		return codes.FailedPrecondition
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	}

	return codes.Unknown
}

var isGrpcAcceptableCode = map[codes.Code]bool{
	codes.DeadlineExceeded: true,
	codes.Internal:         true,
	codes.Unavailable:      true,
	codes.DataLoss:         true,
}

func IsGrpcAcceptable(err error) bool {
	return isGrpcAcceptableCode[status.Code(err)]
}

// GrpcCodeToHTTP gRPC转HTTP Code
func GrpcCodeToHTTP(statusCode codes.Code) int {
	switch statusCode {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusServiceUnavailable
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ConvertErr2Status 内部转换，为了让err=nil的时候，监控数据里有OK信息
func ConvertErr2Status(err *errorpb.ErrResponse) *status.Status {
	if generic.IsNil(err) {
		return status.New(codes.OK, "OK")
	}

	var st, err1 = status.New(codes.Code(err.Code.Code), err.Code.Reason).WithDetails(err)
	if err1 != nil {
		log.Err(err1).Any("lava-error", err).Msg("failed to convert error to grpc status")
		return status.New(codes.Internal, err1.Error())
	}
	return st
}

// ParseError try to convert an error to *Error.
// It supports wrapped errors.
func ParseError(err error) *errorpb.ErrResponse {
	if err == nil {
		return nil
	}

	var code *errorpb.ErrCode
	if errors.As(err, &code) {
		if code.Reason == "" {
			code.Reason = err.Error()
		}

		return &errorpb.ErrResponse{
			Code:   code,
			Detail: errors.ParseToWrap(err),
		}
	}

	return &errorpb.ErrResponse{
		Code: &errorpb.ErrCode{
			Reason: err.Error(),
			Code:   errorpb.Code_Unknown,
			Name:   "lava.unknown",
		},
		Detail: errors.ParseToWrap(err),
	}
}

func Parse(val interface{}) error {
	if generic.IsNil(val) {
		return nil
	}

	switch v := val.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	case []byte:
		return errors.New(convert.B2S(v))
	default:
		return fmt.Errorf("%v", v)
	}
}
