package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubgo/funk/pretty"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/stack"
)

var _ Error = (*baseErr)(nil)

type baseErr struct {
	err    error
	caller *stack.Frame
}

func (t baseErr) String() string {
	if t.err == nil || isNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	if t.err != nil {
		if _err, ok := t.err.(fmt.Stringer); ok {
			buf.WriteString(_err.String())
		} else {
			buf.WriteString(fmt.Sprintf(" %s]: %q\n", color.Red.P("error"), t.err.Error()))
			buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Red.P("detail"), pretty.Sprint(t.err)))
		}
	}

	if t.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", color.Green.P("caller"), t.caller.String()))
	}
	return buf.String()
}

func (t baseErr) MarshalJSON() ([]byte, error) {
	return jjson.Marshal(t.getData())
}

func (t baseErr) getData() map[string]any {
	var data = make(map[string]any)
	if t.caller != nil {
		data["caller"] = t.caller
	}

	if _err, ok := t.err.(json.Marshaler); ok {
		data["cause"] = _err
	} else {
		data["msg"] = t.err.Error()
		data["detail"] = pretty.Sprint(t.err)
	}

	return data
}

func (t baseErr) Unwrap() error { return t.err }

// Error
func (t baseErr) Error() string {
	if t.err == nil || isNil(t.err) {
		return ""
	}

	return t.err.Error()
}
