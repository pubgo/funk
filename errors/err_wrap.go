package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"

	"github.com/pubgo/funk/v2/internal/errors/errinter"

	"github.com/pubgo/funk/v2/proto/errorpb"
	"github.com/pubgo/funk/v2/stack"
)

var (
	_ Error         = (*ErrWrap)(nil)
	_ fmt.Formatter = (*ErrWrap)(nil)
)

func newErrWrapStack(err error, tags Tags) *ErrWrap {
	if err == nil {
		return nil
	}

	pb := &errorpb.ErrWrap{
		Caller: stack.Caller(2).String(),
		Stacks: lo.Map(getStack(), func(item *stack.Frame, index int) string { return item.String() }),
		Error:  MustProtoToAny(ParseErrToPb(err)),
		Tags:   tags.ToMap(),
		Id:     lo.ToPtr(getErrorId(err)),
	}

	return &ErrWrap{err: err, pb: pb}
}

func newErrWrap(err error, tags Tags, callers ...int) *ErrWrap {
	if err == nil {
		return nil
	}

	pb := &errorpb.ErrWrap{
		Caller: stack.Caller(2 + lo.FirstOrEmpty(callers)).String(),
		Error:  MustProtoToAny(ParseErrToPb(err)),
		Tags:   tags.ToMap(),
		Id:     lo.ToPtr(getErrorId(err)),
	}

	return &ErrWrap{err: handleGrpcError(err), pb: pb}
}

type ErrWrap struct {
	err error
	pb  *errorpb.ErrWrap
}

func (e *ErrWrap) ID() string                    { return e.pb.GetId() }
func (e *ErrWrap) Proto() proto.Message          { return e.pb }
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
