package errors

import (
	"bytes"
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/stack"
)

func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return &ErrWrap{
		caller: stack.Caller(1),
		err: &muErrMsg{
			errs: errs,
		},
	}
}

var (
	_ Error         = (*ErrMsg)(nil)
	_ fmt.Formatter = (*ErrMsg)(nil)
)

type muErrMsg struct {
	errs []error
}

func (t *muErrMsg) Unwraps() []error              { return t.errs }
func (t *muErrMsg) Error() string                 { return fmt.Sprintf("%+v", t.errs) }
func (t *muErrMsg) Kind() string                  { return "err_muti" }
func (t *muErrMsg) Format(f fmt.State, verb rune) { strFormat(f, verb, t) }

func (t *muErrMsg) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, t.Kind()))
	for _, err := range t.errs {
		errStringify(buf, err)
	}
	return buf.String()
}

func (t *muErrMsg) MarshalJSON() ([]byte, error) {
	data := errJsonify(t.err)
	data["kind"] = t.Kind()
	data["msg"] = t.pb.Msg
	data["detail"] = t.pb.Detail
	data["stack"] = t.pb.Stack
	data["tags"] = t.pb.Tags
	return json.Marshal(data)
}
