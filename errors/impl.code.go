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
	"github.com/pubgo/funk/proto/errorpb"
	"github.com/pubgo/funk/stack"
)

var _ ErrCode = (*errCodeImpl)(nil)
var _ fmt.Formatter = (*errCodeImpl)(nil)

type errCodeImpl struct {
	err    error
	caller *stack.Frame

	reason string
	name   string
	status uint32
	code   errorpb.Code
	tags   map[string]string
}

func (t *errCodeImpl) Name() string {
	return t.name
}

func (t *errCodeImpl) Status() uint32 {
	return t.status
}

func (t *errCodeImpl) SetErr(err error) ErrCode {
	t.err = err
	return t
}

func (t *errCodeImpl) Tags() map[string]string {
	return t.tags
}

func (t *errCodeImpl) AddTag(key string, val string) ErrCode {
	if t.tags == nil {
		t.tags = make(map[string]string)
	}
	t.tags[key] = val
	return t
}

func (t *errCodeImpl) Code() errorpb.Code {
	return t.code
}

func (t *errCodeImpl) Reason() string {
	return t.reason
}

func (t *errCodeImpl) SetCode(code errorpb.Code) ErrCode {
	t.code = code
	return t
}

func (t *errCodeImpl) SetReason(reason string) ErrCode {
	t.reason = reason
	return t
}

func (t *errCodeImpl) SetName(name string) ErrCode {
	t.name = name
	return t
}

func (t *errCodeImpl) Format(f fmt.State, verb rune) {
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

func (t *errCodeImpl) String() string {
	if generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	if t.code != 0 {
		buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("code"), t.code.String()))
	}

	if t.name != "" {
		buf.WriteString(fmt.Sprintf("  %s]: %s\n", color.Green.P("name"), t.name))
	}

	if t.reason != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", color.Green.P("reason"), t.reason))
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

func (t *errCodeImpl) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["name"] = t.name
	data["code"] = t.code.String()
	data["reason"] = t.reason
	return jjson.Marshal(data)
}

func (t *errCodeImpl) getData() map[string]any {
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

func (t *errCodeImpl) Unwrap() error { return t.err }

// Error
func (t *errCodeImpl) Error() string {
	if generic.IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}
