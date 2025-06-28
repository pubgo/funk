package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func NewErr(msg *Err) error {
	if generic.IsNil(msg) {
		return nil
	}

	if msg.id == "" {
		msg.id = xid.New().String()
	}

	return msg
}

var (
	_ fmt.Formatter = (*Err)(nil)
	_ Error         = (*Err)(nil)
)

type Err struct {
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
	Tags   Tags   `json:"tags,omitempty"`
	id     string
}

func (e Err) ID() string { return e.id }

func (e Err) Proto() proto.Message {
	return &errorpb.ErrMsg{
		Id:     lo.ToPtr(e.id),
		Msg:    e.Msg,
		Detail: e.Detail,
		Tags:   e.Tags.ToMap(),
	}
}

func (e Err) Kind() string                  { return "simple" }
func (e Err) Error() string                 { return e.Msg }
func (e Err) Format(f fmt.State, verb rune) { strFormat(f, verb, e) }

func (e Err) MarshalJSON() ([]byte, error) {
	data := make(map[string]any, 4)
	data["kind"] = e.Kind()
	data["msg"] = e.Msg
	data["detail"] = e.Detail
	data["tags"] = e.Tags
	data["id"] = e.ID()
	return json.Marshal(data)
}

func (e Err) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, e.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorMsg, e.Msg))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorId, e.ID()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorDetail, e.Detail))
	for i := range e.Tags {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorTags, e.Tags[i].String()))
	}
	return buf.String()
}
