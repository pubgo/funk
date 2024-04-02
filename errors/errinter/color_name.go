package errinter

import (
	"strings"

	"github.com/pubgo/funk/internal/color"
)

var width = 10

func Get(name string) string {
	if width-len(name) < 0 {
		return name
	}
	return strings.Repeat(" ", width-len(name)) + name
}

var (
	ColorKind       = color.Green.P(Get("kind"))
	ColorMsg        = color.Green.P(Get("msg"))
	ColorService    = color.Green.P(Get("service"))
	ColorOperation  = color.Green.P(Get("operation"))
	ColorId         = color.Green.P(Get("id"))
	ColorDetail     = color.Green.P(Get("detail"))
	ColorTags       = color.Green.P(Get("tags"))
	ColorErrMsg     = color.Red.P(Get("err_msg"))
	ColorErrDetail  = color.Red.P(Get("err_detail"))
	ColorCaller     = color.Green.P(Get("caller"))
	ColorCode       = color.Green.P(Get("code"))
	ColorMessage    = color.Green.P(Get("message"))
	ColorBiz        = color.Green.P(Get("biz_code"))
	ColorStatusCode = color.Green.P(Get("status_code"))
	ColorName       = color.Green.P(Get("name"))
	ColorStack      = color.Green.P(Get("stack"))
	ColorVersion    = color.Green.P(Get("version"))
)
