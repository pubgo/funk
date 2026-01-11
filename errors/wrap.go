package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/internal/errors/errcolorfield"
	"github.com/pubgo/funk/v2/stack"
)

var (
	_ Error         = (*ErrWrap)(nil)
	_ fmt.Formatter = (*ErrWrap)(nil)
)

func newErrWrapStack(err error, tags Tags) *ErrWrap {
	if err == nil {
		return nil
	}

	return &ErrWrap{
		Caller: stack.Caller(2).String(),
		Stacks: lo.Map(getStack(), func(item *stack.Frame, index int) string { return item.String() }),
		Err:    err,
		Tags:   tags,
		errId:  GetErrorId(err),
	}
}

func newErrWrap(err error, tags Tags, callers ...int) *ErrWrap {
	if err == nil {
		return nil
	}

	return &ErrWrap{
		Caller: stack.Caller(2 + lo.FirstOrEmpty(callers)).String(),
		Err:    err,
		Tags:   tags,
		errId:  GetErrorId(err),
	}
}

type ErrWrap struct {
	Err    error
	Tags   Tags
	Caller string
	Stacks []string
	errId  string
}

func (e *ErrWrap) ID() string                    { return e.errId }
func (e *ErrWrap) Format(f fmt.State, verb rune) { PrintFormat(f, verb, e) }
func (e *ErrWrap) Unwrap() error                 { return e.Err }
func (e *ErrWrap) Error() string                 { return e.Err.Error() }

func (e *ErrWrap) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("===============================================================\n")
	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorId, e.ID())
	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorCaller, e.Caller)
	for k, v := range e.Tags.ToMapString() {
		fmt.Fprintf(buf, "%s]: %s=%q\n", errcolorfield.ColorTags, k, v)
	}

	for i := range e.Stacks {
		fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorStack, e.Stacks[i])
	}
	ErrStringify(buf, e.Err)
	return buf.String()
}

func (e *ErrWrap) MarshalJSON() ([]byte, error) {
	data := ErrJsonify(e.Err)
	if len(e.Tags) > 0 {
		data["fields"] = e.Tags
	}

	if len(e.Stacks) > 0 {
		data["stacks"] = e.Stacks
	}

	data["caller"] = e.Caller
	data["id"] = e.ID()
	return json.Marshal(data)
}
