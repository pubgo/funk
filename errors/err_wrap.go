package errors

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/stack"
)

var _ Error = (*ErrWrap)(nil)
var _ fmt.Formatter = (*ErrWrap)(nil)

type ErrWrap struct {
	err    error
	caller *stack.Frame
	stack  []*stack.Frame
	fields map[string]any
}

func (e *ErrWrap) Format(f fmt.State, verb rune) {
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

func (e *ErrWrap) Kind() string {
	return "err_wrap"
}

func (e *ErrWrap) Error() string {
	return e.err.Error()
}

func (e *ErrWrap) String() string {
	if e == nil {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, e.Kind()))
	if e.caller != nil {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, e.caller.String()))
	}

	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorTags, repr.String(e.fields)))
	for i := range e.stack {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStack, e.stack[i].String()))
	}

	buf.WriteString("\n")
	errStringify(buf, e.err)
	buf.WriteString("====================================================================\n")
	return buf.String()
}

func (e *ErrWrap) Unwrap() error {
	return e.err
}

func (e *ErrWrap) MarshalJSON() ([]byte, error) {
	var data = e.getData()
	data["fields"] = e.fields
	data["kind"] = e.Kind()
	data["stacks"] = e.stack
	data["caller"] = e.caller.String()
	return jjson.Marshal(data)
}

func (e *ErrWrap) getData() map[string]any {
	var data = make(map[string]any)
	data["kind"] = e.Kind()

	var mm = errJsonify(e.err)
	if mm != nil {
		for k, v := range mm {
			data[k] = v
		}
	}

	return data
}
