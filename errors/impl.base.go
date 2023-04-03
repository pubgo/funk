package errors

import (
	"bytes"
	"fmt"
	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

var _ Error = (*baseErr)(nil)
var _ fmt.Formatter = (*baseErr)(nil)

type baseErr struct {
	err    error
	caller *stack.Frame
	msg    string
}

func (t *baseErr) Kind() string {
	return "base"
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
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))
	if t.msg != "" {
		buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorMsg, t.msg))
	}

	if t.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, t.caller.String()))
	}

	errStringify(buf, t.err)

	return buf.String()
}

func (t *baseErr) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["msg"] = t.msg
	return jjson.Marshal(data)
}

func (t *baseErr) getData() map[string]any {
	var data = make(map[string]any)
	data["kind"] = t.Kind()
	if t.caller != nil {
		data["caller"] = t.caller
	}

	var mm = errJsonify(t.err)
	if mm != nil {
		for k, v := range mm {
			data[k] = v
		}
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
