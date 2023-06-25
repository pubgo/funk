package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/stack"
)

func WrapTag(err error, tags ...Tag) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: tags,
	}
}

func WrapFn(err error, fn func() Tags) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: fn(),
	}
}

func WrapKV(err error, key string, value any) error {
	if generic.IsNil(err) {
		return nil
	}

	return &ErrWrap{
		err:    err,
		caller: stack.Caller(1),
		fields: Tags{T(key, value)},
	}
}

var _ Error = (*ErrWrap)(nil)
var _ fmt.Formatter = (*ErrWrap)(nil)

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
	var buf = bytes.NewBuffer(nil)
	buf.WriteString("===============================================================\n")
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, e.Kind()))
	buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorCaller, e.caller.String()))
	for i := range e.fields {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorTags, e.fields[i].String()))
	}

	for i := range e.stack {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStack, e.stack[i].String()))
	}
	errStringify(buf, e.err)
	return buf.String()
}

func (e *ErrWrap) MarshalJSON() ([]byte, error) {
	var data = errJsonify(e.err)
	data["kind"] = e.Kind()
	data["fields"] = e.fields
	data["stacks"] = e.stack
	data["caller"] = e.caller.String()
	return jjson.Marshal(data)
}
