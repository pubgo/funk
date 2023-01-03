package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/internal/color"
)

var _ ITagWrap = (*errTagImpl)(nil)

type errTagImpl struct {
	*errImpl
	tags map[string]any
}

func (e errTagImpl) Tags() map[string]any {
	return e.tags
}

func (e errTagImpl) MarshalJSON() ([]byte, error) {
	var data = e.getData()
	data["tags"] = e.tags
	return jjson.Marshal(data)
}

func (e errTagImpl) String() string {
	if e.err == nil || isNil(e.err) {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("  %s]: %q\n", color.Green.P("tags"), e.tags))
	buf.WriteString(e.errImpl.String())
	return buf.String()
}
