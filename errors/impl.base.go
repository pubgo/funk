package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

var _ XError = (*baseErr)(nil)
var _ fmt.Formatter = (*baseErr)(nil)

type baseErr struct {
	err    error
	caller *stack.Frame

	bizCode string
	code    errorpb.Code
	msg     string
	stacks  []*stack.Frame
	tags    Tags
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
	if t.err == nil || IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	if t.code != 0 {
		buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("code"), t.code.String()))
	}

	if t.bizCode != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("biz"), t.bizCode))
	}

	if t.msg != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %q\n", color.Green.P("msg"), t.msg))
	}

	if t.tags != nil && len(t.tags) > 0 {
		buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("tags"), pretty.Sprint(t.tags)))
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

	for i := range t.stacks {
		buf.WriteString(fmt.Sprintf(" %s]: %s\n", color.Yellow.P("stack"), t.stacks[i].String()))
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
	data["biz_code"] = t.bizCode
	data["code"] = t.code.String()
	data["msg"] = t.msg
	data["stacks"] = t.stacks
	data["tags"] = t.tags
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
	if t.err == nil || IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}
