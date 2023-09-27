package errors

import (
	"bytes"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

func NewRedirectErr(err *errorpb.ErrRedirect) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err:    &ErrRedirect{pb: err},
	}
}

var _ Error = (*ErrRedirect)(nil)
var _ fmt.Formatter = (*ErrRedirect)(nil)

type ErrRedirect struct {
	err error
	pb  *errorpb.ErrRedirect
}

func (t *ErrRedirect) Unwrap() error                 { return t.err }
func (t *ErrRedirect) Error() string                 { return t.err.Error() }
func (t *ErrRedirect) Proto() *errorpb.ErrRedirect   { return t.pb }
func (t *ErrRedirect) Kind() string                  { return "err_redirect" }
func (t *ErrRedirect) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *ErrRedirect) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(t.pb.String())
	errStringify(buf, t.err)
	return buf.String()
}

func (t *ErrRedirect) MarshalJSON() ([]byte, error) {
	var data = errJsonify(t.err)
	data["kind"] = t.Kind()
	return json.Marshal(data)
}
