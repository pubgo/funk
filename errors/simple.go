package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rs/xid"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/internal/errors/errinter"
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
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorId, e.id))
	for k, v := range e.Tags.ToMapString() {
		buf.WriteString(fmt.Sprintf("%s]: %s: %q\n", errinter.ColorTags, k, v))
	}
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorErrMsg, e.Msg))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorErrDetail, e.Detail))
	return buf.String()
}
