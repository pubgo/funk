package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rs/xid"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/internal/errors/errcolorfield"
)

var (
	_ fmt.Formatter = (*Err)(nil)
	_ Error         = (*Err)(nil)
)

func newSimpleErr(err *Err) *Err {
	if funk.IsNil(err) {
		return nil
	}

	if err.Msg == "" {
		err.Msg = err.Detail
	}

	if err.Msg == "" {
		log.Panicln("error message cannot be empty")
	}

	if err.id == "" {
		err.id = xid.New().String()
	}

	return err
}

type Err struct {
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
	Tags   Tags   `json:"tags,omitempty"`
	id     string
}

func (e Err) ID() string                    { return e.id }
func (e Err) Error() string                 { return e.Msg }
func (e Err) Format(f fmt.State, verb rune) { PrintFormat(f, verb, e) }

func (e Err) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"msg":    e.Msg,
		"detail": e.Detail,
		"tags":   e.Tags,
		"id":     e.id,
	})
}

func (e Err) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorId, e.id)
	for k, v := range e.Tags.ToMapString() {
		fmt.Fprintf(buf, "%s]: %s: %q\n", errcolorfield.ColorTags, k, v)
	}
	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorErrMsg, e.Msg)
	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorErrDetail, e.Detail)
	return buf.String()
}
