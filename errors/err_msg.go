package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

func NewMsgErr(msg *errorpb.ErrMsg) error {
	if generic.IsNil(msg) {
		return nil
	}

	if msg.Id == nil {
		msg.Id = newErrorId()
	}

	return &ErrMsg{pb: msg, err: errors.New(msg.Msg)}
}

func WrapMsg(err error, msg *errorpb.ErrMsg) error {
	if err == nil {
		return nil
	}

	if msg == nil {
		return Wrap(err, "error msg is nil")
	}

	if msg.Id == nil {
		msg.Id = lo.ToPtr(getErrorId(err))
	}

	return &ErrWrap{
		err: &ErrMsg{pb: msg, err: handleGrpcError(err)},
		pb: &errorpb.ErrWrap{
			Caller: stack.Caller(1).String(),
			Error:  MustProtoToAny(msg),
		},
	}
}

var (
	_ Error         = (*ErrMsg)(nil)
	_ fmt.Formatter = (*ErrMsg)(nil)
)

type ErrMsg struct {
	err error
	pb  *errorpb.ErrMsg
}

func (t *ErrMsg) ID() string                    { return lo.FromPtr(t.pb.Id) }
func (t *ErrMsg) Unwrap() error                 { return t.err }
func (t *ErrMsg) Error() string                 { return t.err.Error() }
func (t *ErrMsg) Kind() string                  { return "err_msg" }
func (t *ErrMsg) Proto() proto.Message          { return t.pb }
func (t *ErrMsg) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrMsg) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorMsg, t.pb.Msg))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorDetail, t.pb.Detail))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorStack, t.pb.Stack))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorTags, t.pb.Tags))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorId, lo.FromPtr(t.pb.Id)))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrMsg) MarshalJSON() ([]byte, error) {
	data := errJsonify(t.err)
	data["kind"] = t.Kind()
	data["msg"] = t.pb.Msg
	data["detail"] = t.pb.Detail
	data["stack"] = t.pb.Stack
	data["tags"] = t.pb.Tags
	data["id"] = t.pb.Id
	return json.Marshal(data)
}
