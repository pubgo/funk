package strutil

import (
	"io"

	"github.com/pubgo/funk/convert"
	"github.com/valyala/fasttemplate"
)

func Format(template string, data map[string]string) string {
	tpl := fasttemplate.New(template, "{", "}")
	return tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		return w.Write(convert.S2B(data[tag]))
	})
}
