package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func NewCodeErrWithMap(code *errorpb.ErrCode, details ...map[string]any) error {
	code = cloneAndCheck(code)
	if code == nil {
		return nil
	}

	var detailMaps = make(map[string]any)
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
	if generic.IsNil(code) {
		return nil
	}

	if code.Id == nil {
		code.Id = newErrorId()
	}

	if len(details) > 0 {
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
	if err == nil {
		return nil
	}

	code = cloneAndCheck(code)
	if code == nil {
		return Wrap(err, "error code is nil")
	}

	if code.Id == nil {
		code.Id = lo.ToPtr(getErrorId(err))
	}

	code.Details = append(code.Details, MustProtoToAny(ParseErrToPb(err)))
	return &ErrWrap{
		err: &ErrCode{pb: code, err: errors.New(code.Message)},
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(code),
			Id:     code.Id,
		},
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

func (t *ErrCode) ID() string                    { return lo.FromPtr(t.pb.Id) }
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

	var check = func(err2 *ErrCode) bool {
		return err2.pb.Code == t.pb.Code && err2.pb.Name == t.pb.Name
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
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %d\n", errinter.ColorCode, t.pb.Code))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorMessage, t.pb.Message))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorName, t.pb.Name))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorStatusCode, t.pb.StatusCode.String()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorId, lo.FromPtr(t.pb.Id)))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	data := errJsonify(t.err)
	data["kind"] = t.Kind()
	data["name"] = t.pb.Name
	data["status_code"] = t.pb.StatusCode.String()
	data["code"] = t.pb.Code
	data["id"] = t.pb.Id
	data["message"] = t.pb.Message
	return json.Marshal(data)
}
