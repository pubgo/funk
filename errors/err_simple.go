package errors

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
)

var _ fmt.Formatter = (*Err)(nil)

type Err struct {
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
	Tags   Tags   `json:"tags,omitempty"`
}

func (e Err) Kind() string  { return "simple" }
func (e Err) Error() string { return e.Msg }

func (e Err) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = e.MarshalJSON()
		if err != nil {
			fmt.Fprintln(f, err.Error())
		} else {
			fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		fmt.Fprintln(f, e.String())
	}
}

func (e Err) MarshalJSON() ([]byte, error) {
	var data = make(map[string]any, 4)
	data["kind"] = e.Kind()
	data["msg"] = e.Msg
	data["detail"] = e.Detail
	data["tags"] = e.Tags
	return jjson.Marshal(data)
}

func (e Err) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, e.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, e.Msg))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorDetail, e.Detail))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorTags, repr.String(e.Tags)))
	return buf.String()
}
