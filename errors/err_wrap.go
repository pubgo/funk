package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"
	"google.golang.org/protobuf/proto"

	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/proto/errorpb"
)

var (
	_ Error         = (*ErrWrap)(nil)
	_ fmt.Formatter = (*ErrWrap)(nil)
)

type ErrWrap struct {
	err error
	pb  *errorpb.ErrWrap
}

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
	data["fields"] = e.pb.Tags
	data["stacks"] = e.pb.Stacks
	data["caller"] = e.pb.Caller
	return jjson.Marshal(data)
}
