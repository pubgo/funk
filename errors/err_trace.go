package errors

import (
	"bytes"
	"errors"
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
	"google.golang.org/protobuf/proto"
)

func NewTraceErr(trace *errorpb.ErrTrace) error {
	if generic.IsNil(trace) {
		return nil
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrTrace{pb: trace, err: errors.New(trace.String())},
	}
}

func WrapTrace(err error, trace *errorpb.ErrTrace) error {
	if generic.IsNil(err) {
		return nil
	}

	if trace == nil {
		panic("error trace is nil")
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrTrace{pb: trace, err: handleGrpcError(err)},
	}
}

var (
	_ Error         = (*ErrTrace)(nil)
	_ fmt.Formatter = (*ErrTrace)(nil)
)

type ErrTrace struct {
	err error
	pb  *errorpb.ErrTrace
}

func (t *ErrTrace) Unwrap() error                 { return t.err }
func (t *ErrTrace) Error() string                 { return t.err.Error() }
func (t *ErrTrace) Kind() string                  { return "err_trace" }
func (t *ErrTrace) Proto() proto.Message          { return t.pb }
func (t *ErrTrace) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrTrace) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, t.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorId, t.pb.Id))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorOperation, t.pb.Operation))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorService, t.pb.Service))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorVersion, t.pb.Version))
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrTrace) MarshalJSON() ([]byte, error) {
	data := errJsonify(t.err)
	data["kind"] = t.Kind()
	data["id"] = t.pb.Id
	data["operation"] = t.pb.Operation
	data["service"] = t.pb.Service
	data["version"] = t.pb.Version
	return json.Marshal(data)
}
