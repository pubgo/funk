package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

var (
	_ Error         = (*ErrWrap)(nil)
	_ fmt.Formatter = (*ErrWrap)(nil)
)

type ErrWrap struct {
	err error
	pb  *errorpb.ErrWrap
}

func (e *ErrWrap) ID() string { return lo.FromPtr(e.pb.Id) }
func (e *ErrWrap) Proto() proto.Message {
	return e.pb
}
func (e *ErrWrap) Format(f fmt.State, verb rune) { strFormat(f, verb, e) }
func (e *ErrWrap) Unwrap() error                 { return e.err }
func (e *ErrWrap) Kind() string                  { return "err_wrap" }
func (e *ErrWrap) Error() string                 { return e.err.Error() }

func (e *ErrWrap) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("===============================================================\n")
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, e.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorId, e.ID()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorCaller, e.pb.Caller))
	for k, v := range e.pb.Tags {
		buf.WriteString(fmt.Sprintf("%s]: %s=%q\n", errinter.ColorTags, k, v))
	}

	for i := range e.pb.Stacks {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorStack, e.pb.Stacks[i]))
	}
	errStringify(buf, e.err)
	return buf.String()
}

func (e *ErrWrap) MarshalJSON() ([]byte, error) {
	data := errJsonify(e.err)
	data["kind"] = e.Kind()
	if len(e.pb.Tags) > 0 {
		data["fields"] = e.pb.Tags
	}

	if len(e.pb.Stacks) > 0 {
		data["stacks"] = e.pb.Stacks
	}

	data["caller"] = e.pb.Caller
	data["id"] = e.ID()
	return json.Marshal(data)
}
