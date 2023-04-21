package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/proto/errorpb"
)

var _ Error = (*ErrTrace)(nil)
var _ fmt.Formatter = (*ErrTrace)(nil)

type ErrTrace struct {
	err error
	pb  *errorpb.ErrTrace
}

func (t *ErrTrace) Unwrap() error                 { return t.err }
func (t *ErrTrace) Error() string                 { return t.err.Error() }
func (t *ErrTrace) Kind() string                  { return "err_trace" }
func (t *ErrTrace) Proto() *errorpb.ErrTrace      { return t.pb }
func (t *ErrTrace) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrTrace) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorId, t.pb.Id))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorOperation, t.pb.Operation))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorService, t.pb.Service))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorVersion, t.pb.Version))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrTrace) MarshalJSON() ([]byte, error) {
	var data = errJsonify(t.err)
	data["id"] = t.pb.Id
	data["operation"] = t.pb.Operation
	data["service"] = t.pb.Service
	data["version"] = t.pb.Version
	return jjson.Marshal(data)
}
