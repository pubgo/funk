package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/proto/errorpb"
)

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
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStatus, t.pb.Status))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrCode) MarshalJSON() ([]byte, error) {
	var data = errJsonify(t.err)
	data["kind"] = t.Kind()
	data["status"] = t.pb.Status
	data["code"] = t.pb.Code.String()
	data["reason"] = t.pb.Reason
	return jjson.Marshal(data)
}
