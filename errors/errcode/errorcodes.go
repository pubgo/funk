// Package errcode provides error code handling with gRPC compatibility.
package errcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log/logutil"
	"github.com/pubgo/funk/v2/proto/errorpb"

	"github.com/pubgo/funk/v2/internal/errors/errinter"

	"github.com/pubgo/funk/v2"
)

func MustTagsToAny(tags errors.Tags) []*anypb.Any {
	return lo.MapToSlice(tags, func(key string, value any) *anypb.Any {
		return MustProtoToAny(&errorpb.Tag{Key: key, Value: fmt.Sprintf("%v", value)})
	})
}

func MustStructToAny(p map[string]any) *anypb.Any {
	if p == nil {
		return nil
	}

	pb, err := structpb.NewStruct(p)
	if err != nil {
		log.Err(err).
			Any("params", p).
			Func(logutil.WithNotice()).
			Stack().
			Msgf("failed to encode map-any to struct protobuf")
		return anyToProtobuf(p)
	} else {
		return MustProtoToAny(pb)
	}
}

func anyToProtobuf(v any) *anypb.Any {
	data, err := json.Marshal(v)
	if err != nil {
		return &anypb.Any{
			TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
			Value:   []byte(fmt.Sprintf("err:%s detail:%#v", err.Error(), v)),
		}
	}
	return &anypb.Any{
		TypeUrl: "type.googleapis.com/google.protobuf.Struct",
		Value:   data,
	}
}

func ParseErrToPb(err error) proto.Message {
	switch err1 := err.(type) {
	case nil:
		return nil
	case ErrorProto:
		return err1.Proto()
	case GRPCStatus:
		return err1.GRPCStatus().Proto()
	case proto.Message:
		return err1
	default:
		return &errorpb.ErrMsg{Msg: err.Error(), Detail: fmt.Sprintf("%v", err)}
	}
}

func MustProtoToAny(p proto.Message) *anypb.Any {
	switch p := p.(type) {
	case nil:
		return nil
	case *anypb.Any:
		return p
	}

	pb, err := anypb.New(p)
	if err != nil {
		params := prototext.Format(p)
		log.Err(err).
			Str("protobuf", params).
			Func(logutil.WithNotice()).
			Stack().
			Msgf("failed to encode protobuf message to any protobuf")
		pb, _ = anypb.New(structpb.NewStringValue(params))
		return pb
	}

	return pb
}

func handleGrpcError(err error) error {
	switch v := err.(type) {
	case nil:
		return nil
	case *errors.ErrWrap:
		return v
	case GRPCStatus:
		return NewCodeErr(&errorpb.ErrCode{
			Message:    v.GRPCStatus().Message(),
			StatusCode: errorpb.Code(v.GRPCStatus().Code()),
			Name:       "lava.grpc.status",
			Details:    v.GRPCStatus().Proto().Details,
		})
	default:
		return err
	}
}

// ErrorProto is an interface for errors that can be converted to protobuf messages.
type ErrorProto interface {
	error
	Proto() proto.Message
}

// GRPCStatus is an interface for errors that have gRPC status.
type GRPCStatus interface {
	GRPCStatus() *status.Status
}

func NewCodeErrWithMap(code *errorpb.ErrCode, details ...map[string]any) error {
	code = cloneAndCheck(code)
	if code == nil {
		return nil
	}

	detailMaps := make(map[string]any)
	for _, detail := range details {
		for k, v := range detail {
			detailMaps[k] = v
		}
	}

	data, err := structpb.NewStruct(detailMaps)
	if err != nil {
		return NewCodeErr(code,
			structpb.NewStringValue(err.Error()),
			structpb.NewStringValue(fmt.Sprintf("%v", detailMaps)))
	}

	return NewCodeErr(code, data)
}

func NewCodeErrWithMsg(code *errorpb.ErrCode, msg string, details ...proto.Message) error {
	code = cloneAndCheck(code)
	if code == nil {
		return nil
	}

	code.Message = strings.ToTitle(strings.TrimSpace(msg))
	return NewCodeErr(code, details...)
}

func NewCodeErr(code *errorpb.ErrCode, details ...proto.Message) error {
	code = cloneAndCheck(code)
	if funk.IsNil(code) {
		return nil
	}

	if code.Id == nil {
		code.Id = lo.ToPtr(errors.NewErrorId())
	}

	if len(details) > 0 {
		for _, p := range details {
			if p == nil || funk.IsNil(p) {
				continue
			}

			pb := MustProtoToAny(p)
			if pb == nil {
				continue
			}
			code.Details = append(code.Details, pb)
		}
	}

	return &ErrCode{pb: code, err: errors.New(code.Message)}
}

func WrapCode(err error, code *errorpb.ErrCode) error {
	if err == nil {
		return nil
	}

	code = cloneAndCheck(code)
	if code == nil {
		return errors.Wrap(err, "error code is nil")
	}

	if code.Id == nil {
		code.Id = lo.ToPtr(errors.GetErrorId(err))
	}

	code.Details = append(code.Details, MustProtoToAny(ParseErrToPb(err)))
	return errors.WrapTags(
		&ErrCode{pb: code, err: errors.New(code.Message)},
		errors.Tags{"msg": err.Error()},
	)
}

var (
	_ errors.Error  = (*ErrCode)(nil)
	_ fmt.Formatter = (*ErrCode)(nil)
)

// ErrCode represents an error with code information.
type ErrCode struct {
	err error
	pb  *errorpb.ErrCode
}

func (t *ErrCode) ID() string                    { return lo.FromPtr(t.pb.Id) }
func (t *ErrCode) Unwrap() error                 { return t.err }
func (t *ErrCode) Error() string                 { return t.err.Error() }
func (t *ErrCode) Proto() proto.Message          { return t.pb }
func (t *ErrCode) Format(f fmt.State, verb rune) { errors.PrintFormat(f, verb, t) }

func (t *ErrCode) Is(err error) bool {
	if err == nil {
		return false
	}

	if t.err == err {
		return true
	}

	check := func(errCode *ErrCode) bool {
		return errCode.pb.Code == t.pb.Code && errCode.pb.Name == t.pb.Name
	}
	if err1, ok := err.(*ErrCode); ok && check(err1) {
		return true
	}

	return false
}

func (t *ErrCode) As(err any) bool {
	if err == nil {
		return false
	}

	if err1, ok := err.(*ErrCode); ok { //nolint
		err1.pb = t.pb
		return true
	}

	if err1, ok := err.(**errorpb.ErrCode); ok {
		*err1 = t.pb
		return true
	}

	if err1, ok := err.(*errorpb.ErrCode); ok {
		*err1 = lo.FromPtr(proto.Clone(t.pb).(*errorpb.ErrCode))
		return true
	}

	return false
}

func (t *ErrCode) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %d\n", errinter.ColorCode, t.pb.Code))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorMessage, t.pb.Message))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorName, t.pb.Name))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorStatusCode, t.pb.StatusCode.String()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorId, lo.FromPtr(t.pb.Id)))
	errors.ErrStringify(buf, t.err)
	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	data := errors.ErrJsonify(t.err)
	data["name"] = t.pb.Name
	data["status_code"] = t.pb.StatusCode.String()
	data["code"] = t.pb.Code
	data["id"] = t.pb.Id
	data["message"] = t.pb.Message
	return json.Marshal(data)
}

func cloneAndCheck(code *errorpb.ErrCode) *errorpb.ErrCode {
	if code == nil {
		return nil
	}

	code = proto.Clone(code).(*errorpb.ErrCode)
	if code.Code == 0 {
		code.StatusCode = errorpb.Code_OK
	} else if code.StatusCode == errorpb.Code_OK {
		code.StatusCode = errorpb.Code_Internal
	}

	return code
}

// CloneAndCheck clones and validates an error code.
func CloneAndCheck(code *errorpb.ErrCode) *errorpb.ErrCode {
	return cloneAndCheck(code)
}

// Err2GrpcCode converts a standard Go error into its canonical code. Note that
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

// Http2GrpcCode converts an HTTP status code to a gRPC code.
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

// IsGrpcAcceptable checks if a gRPC error is acceptable for retry.
func IsGrpcAcceptable(err error) bool {
	return isGrpcAcceptableCode[status.Code(err)]
}

// GrpcCodeToHTTP converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func GrpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
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
		log.Warn().Msgf("Unknown gRPC error code: %v", code)
		return http.StatusInternalServerError
	}
}

// ConvertErr2Status converts an error code to a gRPC status.
// Internal conversion, so that when err=nil, OK information is included in monitoring data.
func ConvertErr2Status(errCode *errorpb.ErrCode) *status.Status {
	if funk.IsNil(errCode) {
		return status.New(codes.OK, "OK")
	}

	if (errCode.Name != "" || errCode.Code != 0) && errCode.StatusCode == 0 {
		errCode.StatusCode = errorpb.Code_Internal
	}

	st := status.New(codes.Code(errCode.StatusCode), errCode.Message)
	if st1, err1 := st.WithDetails(errCode); err1 != nil {
		log.Err(err1).Any("lava-error", errCode).Msg("failed to convert error to grpc status")
		return st
	} else {
		return st1
	}
}

// ParseError tries to convert an error to *Error.
// It supports wrapped errors.
func ParseError(err error) *errorpb.ErrCode {
	if err == nil {
		return nil
	}

	// grpc error
	gs, ok := err.(GRPCStatus)
	if ok {
		if gs.GRPCStatus().Code() == codes.OK {
			return nil
		}

		details := gs.GRPCStatus().Details()
		if len(details) > 0 && details[0] != nil {
			if e, ok := details[0].(*errorpb.ErrCode); ok && e != nil {
				return e
			}
		}

		return &errorpb.ErrCode{
			Message:    gs.GRPCStatus().Message(),
			StatusCode: errorpb.Code(gs.GRPCStatus().Code()),
			Code:       int32(GrpcCodeToHTTP(gs.GRPCStatus().Code())),
			Name:       "lava.grpc.status",
			Details:    gs.GRPCStatus().Proto().Details,
		}
	}

	var ce *ErrCode
	if errors.As(err, &ce) {
		pb := ce.Proto().(*errorpb.ErrCode)
		if pb.Message == "" {
			pb.Message = err.Error()
		}

		return pb
	}

	return &errorpb.ErrCode{
		Message:    err.Error(),
		StatusCode: errorpb.Code_Unknown,
		Code:       500,
		Name:       "lava.error.unknown",
		Details:    []*anypb.Any{MustProtoToAny(ParseErrToPb(err))},
	}
}
