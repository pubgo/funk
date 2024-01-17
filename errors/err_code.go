package errors

import (
	"bytes"
	"errors"
	"fmt"
	json "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
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

			pb, err := anypb.New(p)
			if err != nil {
				log.Printf("failed to encode to any")
				continue
			}

			code.Details = append(code.Details, pb)
		}
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrCode{pb: code, err: errors.New(code.Message)},
	}
}

func WrapCode(err error, code *errorpb.ErrCode) error {
	if generic.IsNil(err) {
		return nil
	}

	if code == nil {
		panic("error code is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrCode{pb: code, err: handleGrpcError(err)},
	}
}

var _ Error = (*ErrCode)(nil)
var _ fmt.Formatter = (*ErrCode)(nil)

type ErrCode struct {
	err error
	pb  *errorpb.ErrCode
}

func (t *ErrCode) Unwrap() error                 { return t.err }
func (t *ErrCode) Error() string                 { return t.err.Error() }
func (t *ErrCode) Proto() *errorpb.ErrCode       { return t.pb }
func (t *ErrCode) Kind() string                  { return "err_code" }
func (t *ErrCode) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrCode) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %d\n", internal.ColorCode, t.pb.Code))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMessage, t.pb.Message))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorName, t.pb.Name))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStatusCode, t.pb.StatusCode.String()))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	var data = errJsonify(t.err)
	data["kind"] = t.Kind()
	data["name"] = t.pb.Name
	data["status_code"] = t.pb.StatusCode.String()
	data["code"] = t.pb.Code
	data["message"] = t.pb.Message
	return json.Marshal(data)
}
