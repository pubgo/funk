package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
	"github.com/pubgo/funk/stack"
)

var _ IStackWrap = (*errStackImpl)(nil)

type errStackImpl struct {
	*baseErr
	stacks []*stack.Frame
}

func (e errStackImpl) Stack() []*stack.Frame {
	return e.stacks
}

func (e errStackImpl) MarshalJSON() ([]byte, error) {
	var data = e.getData()
	data["stacks"] = e.stacks
	return jjson.Marshal(data)
}

func (e errStackImpl) String() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(e.baseErr.String())
	for i := range e.stacks {
		buf.WriteString(fmt.Sprintf(" %s]: %s\n", color.Yellow.P("stack"), e.stacks[i].String()))
	}
	return buf.String()
}
