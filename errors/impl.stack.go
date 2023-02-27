package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubgo/funk/errors/internal"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
)

var _ ErrStack = (*errStackImpl)(nil)
var _ fmt.Formatter = (*errStackImpl)(nil)

type errStackImpl struct {
	err    error
	stacks []*stack.Frame
}

func (t *errStackImpl) Kind() string {
	return "stack"
}

func (t *errStackImpl) Format(f fmt.State, verb rune) {
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

func (t *errStackImpl) String() string {
	if generic.IsNil(t.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, t.Kind()))

	if t.err != nil {
		if _, ok := t.err.(fmt.Stringer); !ok {
			buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorErrMsg, t.err.Error()))
			buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrDetail, pretty.Sprint(t.err)))
		}
	}

	for i := range t.stacks {
		buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorStack, t.stacks[i].String()))
	}

	if t.err != nil {
		if _err, ok := t.err.(fmt.Stringer); ok {
			buf.WriteString("====================================================================\n")
			buf.WriteString(_err.String())
		}
	}

	return buf.String()
}

func (t *errStackImpl) MarshalJSON() ([]byte, error) {
	var data = t.getData()
	data["stacks"] = t.stacks
	return jjson.Marshal(data)
}

func (t *errStackImpl) getData() map[string]any {
	var data = make(map[string]any)
	if _err, ok := t.err.(json.Marshaler); ok {
		data["cause"] = _err
	} else if t.err != nil {
		data["err_msg"] = t.err.Error()
		data["err_detail"] = repr.String(t.err)
	}

	return data
}

func (t *errStackImpl) Unwrap() error { return t.err }

// Error
func (t *errStackImpl) Error() string {
	if generic.IsNil(t.err) {
		return ""
	}

	return t.err.Error()
}

func (t *errStackImpl) Stack() []*stack.Frame {
	return t.stacks
}

func (t *errStackImpl) AddStack() {
	for i := 0; ; i++ {
		var cc = stack.Caller(1 + i)
		if cc == nil {
			break
		}

		if cc.IsRuntime() {
			continue
		}

		t.stacks = append(t.stacks, cc)
	}
}
