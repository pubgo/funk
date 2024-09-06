package errors

import (
	"bytes"
	"errors"
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"google.golang.org/protobuf/proto"
)

func NewCodeErr(code *errorpb.ErrCode, details ...proto.Message) error {
	if generic.IsNil(code) {
		return nil
	}

	if len(details) > 0 {
		code = proto.Clone(code).(*errorpb.ErrCode)
		for _, p := range details {
			if p == nil || generic.IsNil(p) {
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
	if generic.IsNil(err) {
		return nil
	}

	if code == nil {
		panic("error code is nil")
	}

	code.Details = append(code.Details, MustProtoToAny(ParseErrToPb(err)))
	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrCode{pb: code, err: handleGrpcError(err)},
	}
}

var (
	_ Error         = (*ErrCode)(nil)
	_ fmt.Formatter = (*ErrCode)(nil)
)

type ErrCode struct {
	err error
	pb  *errorpb.ErrCode
}

func (t *ErrCode) Unwrap() error                 { return t.err }
func (t *ErrCode) Error() string                 { return t.err.Error() }
func (t *ErrCode) Proto() proto.Message          { return t.pb }
func (t *ErrCode) Kind() string                  { return "err_code" }
func (t *ErrCode) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrCode) Is(err error) bool {
	if err == nil {
		return false
	}

	if t.err == err {
		return true
	}

	err1, ok := err.(*ErrCode)
	if ok && err1.pb.Code == t.pb.Code && err1.pb.Name == t.pb.Name {
		return true
	}

	return false
}

func (t *ErrCode) As(err any) bool {
	if err == nil {
		return false
	}

	err1, ok := err.(*ErrCode)
	if ok && err1.pb != nil && err1.pb.Code == t.pb.Code && err1.pb.Name == t.pb.Name {
		return true
	}

	err2, ok := err.(**errorpb.ErrCode)
	if ok && (*err2).Code == t.pb.Code && (*err2).Name == t.pb.Name {
		return true
	}

	if err2, ok := err.(*errorpb.ErrCode); ok && (*err2).Code == t.pb.Code && (*err2).Name == t.pb.Name {
		return true
	}

	return false
}

func (t *ErrCode) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %d\n", errinter.ColorCode, t.pb.Code))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorMessage, t.pb.Message))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorName, t.pb.Name))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorStatusCode, t.pb.StatusCode.String()))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	data := errJsonify(t.err)
	data["kind"] = t.Kind()
	data["name"] = t.pb.Name
	data["status_code"] = t.pb.StatusCode.String()
	data["code"] = t.pb.Code
	data["message"] = t.pb.Message
	return json.Marshal(data)
}
