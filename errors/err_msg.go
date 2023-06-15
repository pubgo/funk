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

func WrapMsg(err error, msg *errorpb.ErrMsg) error {
	if generic.IsNil(err) {
		return nil
	}

	if msg == nil {
		panic("error msg is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrMsg{pb: msg, err: err},
	}
}

var _ Error = (*ErrMsg)(nil)
var _ fmt.Formatter = (*ErrMsg)(nil)

type ErrMsg struct {
	err error
	pb  *errorpb.ErrMsg
}

func (t *ErrMsg) Unwrap() error                 { return t.err }
func (t *ErrMsg) Error() string                 { return t.err.Error() }
func (t *ErrMsg) Kind() string                  { return "err_msg" }
func (t *ErrMsg) Proto() *errorpb.ErrMsg        { return t.pb }
func (t *ErrMsg) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrMsg) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, t.pb.Msg))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorDetail, t.pb.Detail))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorStack, t.pb.Stack))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorTags, t.pb.Tags))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrMsg) MarshalJSON() ([]byte, error) {
	var data = errJsonify(t.err)
	data["kind"] = t.Kind()
	data["msg"] = t.pb.Msg
	data["detail"] = t.pb.Detail
	data["stack"] = t.pb.Stack
	data["tags"] = t.pb.Tags
	return json.Marshal(data)
}
