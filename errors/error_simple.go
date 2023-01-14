package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubgo/funk/internal/color"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
)

func ErrOf(fn func(err *Err)) *Err {
	var err = &Err{Caller: stack.Caller(1)}
	fn(err)
	return err
}

var _ fmt.Formatter = (*Err)(nil)
var _ Error = (*Err)(nil)

type Err struct {
	Caller *stack.Frame `json:"caller,omitempty"`
	Err    error        `json:"err,omitempty"`
	Msg    string       `json:"msg,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Tags   Tags         `json:"tags,omitempty"`
}

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

func (e Err) Unwrap() error {
	return e.Err
}

func (e Err) MarshalJSON() ([]byte, error) {
	var data = make(map[string]any)
	if e.Msg != "" {
		data["msg"] = e.Msg
	}

	if e.Detail != "" {
		data["detail"] = e.Detail
	}

	if e.Caller != nil {
		data["caller"] = e.Caller
	}

	if e.Tags != nil {
		data["tags"] = e.Tags
	}

	if _err, ok := e.Err.(json.Marshaler); ok && _err != nil {
		data["cause"] = _err
	} else if e.Err != nil {
		data["err_msg"] = e.Err.Error()
		data["err_detail"] = repr.String(e.Err)
	}
	return jjson.Marshal(data)
}

func (e Err) Error() string {
	if e.Err == nil {
		return e.Msg
	}

	return e.Err.Error()
}

func (e Err) String() string {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("   %s]: %q\n", color.Green.P("msg"), e.Msg))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("detail"), pretty.Sprint(e.Detail)))
	if e.Tags != nil && len(e.Tags) > 0 {
		buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("tags"), pretty.Sprint(e.Tags)))
	}

	if e.Err != nil {
		if _err, ok := e.Err.(fmt.Stringer); !ok {
			buf.WriteString(fmt.Sprintf(" %s]: %q\n", color.Red.P("err_msg"), e.Err.Error()))
			buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Red.P("err_stack"), pretty.Sprint(e.Err)))
			if e.Caller != nil {
				buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("caller"), e.Caller.String()))
			}
		} else {
			buf.WriteString("\n====================================================================\n")
			buf.WriteString(_err.String())
		}
	}

	return buf.String()
}
