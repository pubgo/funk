package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/pretty"
)

var _ Error = (*Err)(nil)

type Err struct {
	Err    error
	Msg    string
	Detail string
}

func (e Err) Unwrap() error {
	return e.Err
}

func (e Err) MarshalJSON() ([]byte, error) {
	var data = make(map[string]any)
	data["msg"] = e.Msg
	data["detail"] = e.Detail

	if _err, ok := e.Err.(json.Marshaler); ok && _err != nil {
		data["cause"] = _err
	} else if e.Err != nil {
		data["err_msg"] = e.Err.Error()
		data["err_detail"] = fmt.Sprintf("%#v", e.Err)
	}
	return jjson.Marshal(data)
}

func (e Err) Error() string {
	if e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

func (e Err) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("msg=%q detail=%q", e.Msg, e.Detail))
	if e.Err != nil {
		if _err, ok := e.Err.(fmt.Stringer); !ok {
			buf.WriteString(fmt.Sprintf(" err_msg=%q", e.Err.Error()))
			buf.WriteString(fmt.Sprintf(" err_detail=%s", pretty.Sprint(e.Err)))
		} else {
			buf.WriteString("\n====================================================================\n")
			buf.WriteString(_err.String())
		}
	}
	return buf.String()
}
