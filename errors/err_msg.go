package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
)

var _ Error = (*ErrMsg)(nil)
var _ fmt.Formatter = (*ErrMsg)(nil)

type ErrMsg struct {
	*ErrBase
	pb *errorpb.ErrMsg
}

func (t *ErrMsg) Kind() string {
	return "err_msg"
}

func (t *ErrMsg) Format(f fmt.State, verb rune) {
	Format(f, verb, t)
}

func (t *ErrMsg) String() string {
	if generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))
	if t.pb.Msg != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, t.pb.Msg))
	}

	if t.pb.Detail != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, t.pb.Detail))
	}

	if t.pb.Stack != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, t.pb.Stack))
	}

	if t.pb.Tags != nil {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, t.pb.Tags))
	}

	errStringify(buf, t.err)

	return buf.String()
}

func (t *ErrMsg) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["msg"] = t.pb.Msg
	data["detail"] = t.pb.Detail
	data["stack"] = t.pb.Stack
	data["tags"] = t.pb.Tags
	return jjson.Marshal(data)
}
