package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/errinter"
	"github.com/pubgo/funk/stack"
)

var (
	_ Error         = (*ErrWrap)(nil)
	_ fmt.Formatter = (*ErrWrap)(nil)
)

type ErrWrap struct {
	err    error
	caller *stack.Frame
	stack  []*stack.Frame
	fields Tags
}

func (e *ErrWrap) Format(f fmt.State, verb rune) { strFormat(f, verb, e) }
func (e *ErrWrap) Unwrap() error                 { return e.err }
func (e *ErrWrap) Kind() string                  { return "err_wrap" }
func (e *ErrWrap) Error() string                 { return e.err.Error() }

func (e *ErrWrap) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("===============================================================\n")
	buf.WriteString(fmt.Sprintf("%s]: %q\n", errinter.ColorKind, e.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorCaller, e.caller.String()))
	for i := range e.fields {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorTags, e.fields[i].String()))
	}

	for i := range e.stack {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", errinter.ColorStack, e.stack[i].String()))
	}
	errStringify(buf, e.err)
	return buf.String()
}

func (e *ErrWrap) MarshalJSON() ([]byte, error) {
	data := errJsonify(e.err)
	data["kind"] = e.Kind()
	data["fields"] = e.fields
	data["stacks"] = e.stack
	data["caller"] = e.caller.String()
	return jjson.Marshal(data)
}
