package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
)

var _ Error = (*baseErr)(nil)
var _ fmt.Formatter = (*baseErr)(nil)

type baseErr struct {
	err    error
	caller *stack.Frame
	msg    string
}

func (t *baseErr) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = t.MarshalJSON()
		if err != nil {
			fmt.Fprintln(f, err.Error())
		} else {
			fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		fmt.Fprintln(f, t.String())
	}
}

func (t *baseErr) String() string {
	if generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	if t.msg != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %q\n", color.Green.P("msg"), t.msg))
	}

	if t.err != nil {
		if _, ok := t.err.(fmt.Stringer); !ok {
			buf.WriteString(fmt.Sprintf(" %s]: %q\n", color.Red.P("error"), t.err.Error()))
			buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Red.P("detail"), pretty.Sprint(t.err)))
		}
	}

	if t.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("caller"), t.caller.String()))
	}

	if t.err != nil {
		if _err, ok := t.err.(fmt.Stringer); ok {
			buf.WriteString("====================================================================\n")
			buf.WriteString(_err.String())
		}
	}

	return buf.String()
}

func (t *baseErr) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["msg"] = t.msg
	return jjson.Marshal(data)
}

func (t *baseErr) getData() map[string]any {
	var data = make(map[string]any)
	if t.caller != nil {
		data["caller"] = t.caller
	}

	if _err, ok := t.err.(json.Marshaler); ok {
		data["cause"] = _err
	} else if t.err != nil {
		data["err_msg"] = t.err.Error()
		data["err_detail"] = repr.String(t.err)
	}

	return data
}

func (t *baseErr) Unwrap() error { return t.err }

// Error
func (t *baseErr) Error() string {
	if generic.IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}
