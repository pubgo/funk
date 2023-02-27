package errutil

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/version"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func parseError(err error) (op string, goType string, desc string, extra map[string]string) {
	extra = make(map[string]string)

	// interfaces
	if _, ok := err.(net.Error); ok {
		if opError, ok := err.(*net.OpError); ok {
			op = opError.Op
			if opError.Source != nil {
				extra["remote_addr"] = opError.Source.String()
			}
			if opError.Addr != nil {
				extra["local_addr"] = opError.Addr.String()
			}
			extra["network"] = opError.Net
			err = opError.Err
		}
		switch actual := err.(type) {
		case *net.AddrError:
			goType = "net.AddrError"
			desc = actual.Err
			extra["addr"] = actual.Addr
		case *net.DNSError:
			goType = "net.DNSError"
			desc = actual.Err
			extra["domain"] = actual.Name
			if actual.Server != "" {
				extra["dns_server"] = actual.Server
			}
		case *net.InvalidAddrError:
			goType = "net.InvalidAddrError"
			desc = actual.Error()
		case *net.ParseError:
			goType = "net.ParseError"
			desc = "invalid " + actual.Type
			extra["text_to_parse"] = actual.Text
		case net.UnknownNetworkError:
			goType = "net.UnknownNetworkError"
			desc = "unknown network"
		case syscall.Errno:
			goType = "syscall.Errno"
			desc = actual.Error()
		case *url.Error:
			goType = "url.Error"
			desc = actual.Err.Error()
			op = actual.Op
		default:
			goType = reflect.TypeOf(err).String()
			desc = err.Error()
		}
		return
	}
	if _, ok := err.(runtime.Error); ok {
		desc = err.Error()
		switch err.(type) {
		case *runtime.TypeAssertionError:
			goType = "runtime.TypeAssertionError"
		default:
			goType = reflect.TypeOf(err).String()
		}
		return
	}

	// structs
	switch actual := err.(type) {
	case *http.ProtocolError:
		desc = actual.ErrorString
		if name, ok := httpProtocolErrors[err]; ok {
			goType = name
		} else {
			goType = "http.ProtocolError"
		}
	case url.EscapeError, *url.EscapeError:
		goType = "url.EscapeError"
		desc = "invalid URL escape"
	case url.InvalidHostError, *url.InvalidHostError:
		goType = "url.InvalidHostError"
		desc = "invalid character in host name"
	case *textproto.Error:
		goType = "textproto.Error"
		desc = actual.Error()
	case textproto.ProtocolError, *textproto.ProtocolError:
		goType = "textproto.ProtocolError"
		desc = actual.Error()

	case tls.RecordHeaderError:
		goType = "tls.RecordHeaderError"
		desc = actual.Msg
		extra["header"] = hex.EncodeToString(actual.RecordHeader[:])
	case x509.CertificateInvalidError:
		goType = "x509.CertificateInvalidError"
		desc = actual.Error()
	case x509.ConstraintViolationError:
		goType = "x509.ConstraintViolationError"
		desc = actual.Error()
	case x509.HostnameError:
		goType = "x509.HostnameError"
		desc = actual.Error()
		extra["host"] = actual.Host
	case x509.InsecureAlgorithmError:
		goType = "x509.InsecureAlgorithmError"
		desc = actual.Error()
	case x509.SystemRootsError:
		goType = "x509.SystemRootsError"
		desc = actual.Error()
	case x509.UnhandledCriticalExtension:
		goType = "x509.UnhandledCriticalExtension"
		desc = actual.Error()
	case x509.UnknownAuthorityError:
		goType = "x509.UnknownAuthorityError"
		desc = actual.Error()
	case hex.InvalidByteError:
		goType = "hex.InvalidByteError"
		desc = "invalid byte"
	case *json.InvalidUTF8Error:
		goType = "json.InvalidUTF8Error"
		desc = "invalid UTF-8 in string"
	case *json.InvalidUnmarshalError:
		goType = "json.InvalidUnmarshalError"
		desc = actual.Error()
	case *json.MarshalerError:
		goType = "json.MarshalerError"
		desc = actual.Error()
	case *json.SyntaxError:
		goType = "json.SyntaxError"
		desc = actual.Error()
	case *json.UnmarshalFieldError:
		goType = "json.UnmarshalFieldError"
		desc = actual.Error()
	case *json.UnmarshalTypeError:
		goType = "json.UnmarshalTypeError"
		desc = actual.Error()
	case *json.UnsupportedTypeError:
		goType = "json.UnsupportedTypeError"
		desc = actual.Error()
	case *json.UnsupportedValueError:
		goType = "json.UnsupportedValueError"
		desc = actual.Error()

	case *os.LinkError:
		goType = "os.LinkError"
		desc = actual.Error()
	case *os.PathError:
		goType = "os.PathError"
		op = actual.Op
		desc = actual.Err.Error()
	case *os.SyscallError:
		goType = "os.SyscallError"
		op = actual.Syscall
		desc = actual.Err.Error()
	case *exec.Error:
		goType = "exec.Error"
		desc = actual.Err.Error()
	case *exec.ExitError:
		goType = "exec.ExitError"
		desc = actual.Error()
		// TODO: limit the length
		extra["stderr"] = string(actual.Stderr)
	case *strconv.NumError:
		goType = "strconv.NumError"
		desc = actual.Err.Error()
		extra["function"] = actual.Func
	case *time.ParseError:
		goType = "time.ParseError"
		desc = actual.Message
	default:
		desc = err.Error()
		if t, ok := miscErrors[err]; ok {
			goType = t
			return
		}
		goType = reflect.TypeOf(err).String()
	}
	return
}

var miscErrors = map[error]string{
	bufio.ErrInvalidUnreadByte: "bufio.ErrInvalidUnreadByte",
	bufio.ErrInvalidUnreadRune: "bufio.ErrInvalidUnreadRune",
	bufio.ErrBufferFull:        "bufio.ErrBufferFull",
	bufio.ErrNegativeCount:     "bufio.ErrNegativeCount",
	bufio.ErrTooLong:           "bufio.ErrTooLong",
	bufio.ErrNegativeAdvance:   "bufio.ErrNegativeAdvance",
	bufio.ErrAdvanceTooFar:     "bufio.ErrAdvanceTooFar",
	bufio.ErrFinalToken:        "bufio.ErrFinalToken",

	http.ErrWriteAfterFlush:    "http.ErrWriteAfterFlush",
	http.ErrBodyNotAllowed:     "http.ErrBodyNotAllowed",
	http.ErrHijacked:           "http.ErrHijacked",
	http.ErrContentLength:      "http.ErrContentLength",
	http.ErrBodyReadAfterClose: "http.ErrBodyReadAfterClose",
	http.ErrHandlerTimeout:     "http.ErrHandlerTimeout",
	http.ErrLineTooLong:        "http.ErrLineTooLong",
	http.ErrMissingFile:        "http.ErrMissingFile",
	http.ErrNoCookie:           "http.ErrNoCookie",
	http.ErrNoLocation:         "http.ErrNoLocation",
	http.ErrSkipAltProtocol:    "http.ErrSkipAltProtocol",

	io.EOF:              "io.EOF",
	io.ErrClosedPipe:    "io.ErrClosedPipe",
	io.ErrNoProgress:    "io.ErrNoProgress",
	io.ErrShortBuffer:   "io.ErrShortBuffer",
	io.ErrShortWrite:    "io.ErrShortWrite",
	io.ErrUnexpectedEOF: "io.ErrUnexpectedEOF",

	os.ErrInvalid:    "os.ErrInvalid",
	os.ErrPermission: "os.ErrPermission",
	os.ErrExist:      "os.ErrExist",
	os.ErrNotExist:   "os.ErrNotExist",

	exec.ErrNotFound: "exec.ErrNotFound",

	x509.ErrUnsupportedAlgorithm: "x509.ErrUnsupportedAlgorithm",
	x509.IncorrectPasswordError:  "x509.IncorrectPasswordError",

	hex.ErrLength: "hex.ErrLength",
}

var httpProtocolErrors = map[error]string{
	http.ErrHeaderTooLong:        "http.ErrHeaderTooLong",
	http.ErrShortBody:            "http.ErrShortBody",
	http.ErrNotSupported:         "http.ErrNotSupported",
	http.ErrUnexpectedTrailer:    "http.ErrUnexpectedTrailer",
	http.ErrMissingContentLength: "http.ErrMissingContentLength",
	http.ErrNotMultipart:         "http.ErrNotMultipart",
	http.ErrMissingBoundary:      "http.ErrMissingBoundary",
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

func IsGrpcAcceptable(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		return false
	default:
		return true
	}
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

	var st, err1 = status.New(codes.Code(err.Code), err.ErrMsg).WithDetails(err)
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

	var ce errors.ErrCode
	if errors.As(err, &ce) {
		return &errorpb.Error{
			Service:   version.Project(),
			Version:   version.Version(),
			Code:      ce.Code(),
			Status:    ce.Status(),
			Name:      ce.Name(),
			Reason:    ce.Reason(),
			ErrMsg:    err.Error(),
			ErrDetail: []byte(fmt.Sprintf("%#v", err)),
			Tags:      ce.Tags(),
		}
	}

	// grpc error
	gs, ok := err.(interface{ GRPCStatus() *status.Status })
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
			ErrMsg:    err.Error(),
			ErrDetail: []byte(fmt.Sprintf("%v", gs.GRPCStatus().Details())),
			Reason:    gs.GRPCStatus().Message(),
			Code:      errorpb.Code(gs.GRPCStatus().Code()),
			Status:    uint32(gs.GRPCStatus().Code()),
			Name:      "lava.grpc.status",
		}
	}

	return &errorpb.Error{
		ErrMsg:    err.Error(),
		ErrDetail: []byte(fmt.Sprintf("%#v", err)),
		Reason:    err.Error(),
		Code:      errorpb.Code_Unknown,
		Status:    uint32(errorpb.Code_Unknown),
		Name:      "lava.unknown",
	}
}
