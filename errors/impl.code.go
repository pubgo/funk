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

	reason  string
	bizCode string
	code    errorpb.Code
	tags    map[string]string
}

func (t *errCodeImpl) Tags() map[string]string {
	return t.tags
}

func (t *errCodeImpl) AddTag(key string, val string) {
	if t.tags == nil {
		t.tags = make(map[string]string)
	}
	t.tags[key] = val
}

func (t *errCodeImpl) Code() errorpb.Code {
	return t.code
}

func (t *errCodeImpl) BizCode() string {
	return t.bizCode
}

func (t *errCodeImpl) Reason() string {
	return t.reason
}

func (t *errCodeImpl) SetCode(code errorpb.Code) {
	t.code = code
}

func (t *errCodeImpl) SetReason(reason string) {
	t.reason = reason
}

func (t *errCodeImpl) SetBizCode(biz string) {
	t.bizCode = biz
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

	if t.bizCode != "" {
		buf.WriteString(fmt.Sprintf("   %s]: %s\n", color.Green.P("biz"), t.bizCode))
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
	data["biz_code"] = t.bizCode
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
