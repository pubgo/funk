package errors

import (
	"bytes"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func WrapCode(err error, code *errorpb.ErrCode) error {
	if generic.IsNil(err) {
		return nil
	}

	if code == nil {
		panic("error code is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrCode{pb: code, err: err},
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
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCode, t.pb.Code.String()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorReason, t.pb.Reason))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorName, t.pb.Name))
	buf.WriteString(fmt.Sprintf("%s]: %d\n", internal.ColorBiz, t.pb.BizCode))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	var data = errJsonify(t.err)
	data["kind"] = t.Kind()
	data["name"] = t.pb.Name
	data["biz_code"] = t.pb.BizCode
	data["code"] = t.pb.Code.String()
	data["reason"] = t.pb.Reason
	return json.Marshal(data)
}
