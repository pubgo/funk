package errutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/version"
)

func Json(err error) []byte {
	if generic.IsNil(err) {
		return nil
	}

	err = errors.Parse(err)
	data, err := jjson.Marshal(err)
	if err != nil {
		log.Err(err).Stack().Str("err_stack", repr.String(err)).Msg("failed to marshal error")
		panic(fmt.Errorf("failed to marshal error, err=%w", err))
	}
	return data
}

func JsonPretty(err error) []byte {
	if generic.IsNil(err) {
		return nil
	}

	err = errors.Parse(err)
	data, err := jjson.MarshalIndent(err, " ", "  ")
	if err != nil {
		log.Err(err).Stack().Str("err_stack", repr.String(err)).Msg("failed to marshal error")
		panic(fmt.Errorf("failed to marshal error, err=%w", err))
	}
	return data
}

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
	switch {
	case err == nil:
		return codes.OK
	case err == io.EOF:
		return codes.OutOfRange
	case errors.Is(err, io.ErrClosedPipe), errors.Is(err, io.ErrNoProgress), errors.Is(err, io.ErrShortBuffer), errors.Is(err, io.ErrShortWrite), errors.Is(err, io.ErrUnexpectedEOF):
		return codes.FailedPrecondition
	case errors.Is(err, os.ErrInvalid):
		return codes.InvalidArgument
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
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
func ConvertErr2Status(err *errorpb.Error) *status.Status {
	if generic.IsNil(err) {
		return status.New(codes.OK, "OK")
	}

	var st, err1 = status.New(codes.Code(err.Code.StatusCode), err.Msg.Msg).WithDetails(err)
	if err1 != nil {
		log.Err(err1).Any("lava-error", err).Msg("failed to convert error to grpc status")
		return status.New(codes.Internal, err1.Error())
	}
	return st
}

// ParseError try to convert an error to *Error.
// It supports wrapped errors.
func ParseError(err error) *errorpb.Error {
	if err == nil {
		return nil
	}

	var ce *errors.ErrCode
	if errors.As(err, &ce) {
		if ce.Proto().Message == "" {
			ce.Proto().Message = err.Error()
		}

		return &errorpb.Error{
			Code: ce.Proto(),
			Trace: &errorpb.ErrTrace{
				Service: version.Project(),
				Version: version.Version(),
			},
			Msg: &errorpb.ErrMsg{
				Msg:    err.Error(),
				Detail: fmt.Sprintf("%#v", err),
			},
		}
	}

	// grpc error
	gs, ok := err.(errors.GRPCStatus)
	if ok {
		if gs.GRPCStatus().Code() == codes.OK {
			return nil
		}

		details := gs.GRPCStatus().Details()
		if len(details) > 0 && details[0] != nil {
			if e, ok := details[0].(*errorpb.Error); ok && e != nil {
				return e
			}
		}

		return &errorpb.Error{
			Code: &errorpb.ErrCode{
				Message:    gs.GRPCStatus().Message(),
				StatusCode: errorpb.Code(gs.GRPCStatus().Code()),
				Code:       int32(GrpcCodeToHTTP(gs.GRPCStatus().Code())),
				Name:       "lava.grpc.status",
			},
			Trace: &errorpb.ErrTrace{
				Service: version.Project(),
				Version: version.Version(),
			},
			Msg: &errorpb.ErrMsg{
				Msg:    err.Error(),
				Detail: fmt.Sprintf("%v", gs.GRPCStatus().Details()),
			},
		}
	}

	return &errorpb.Error{
		Code: &errorpb.ErrCode{
			Message:    err.Error(),
			StatusCode: errorpb.Code_Unknown,
			Code:       500,
			Name:       "lava.error.unknown",
		},
		Trace: &errorpb.ErrTrace{
			Service: version.Project(),
			Version: version.Version(),
		},
		Msg: &errorpb.ErrMsg{
			Msg:    err.Error(),
			Detail: fmt.Sprintf("%#v", err),
		},
	}
}
