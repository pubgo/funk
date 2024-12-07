package errinter

import (
	"strings"

	"github.com/pubgo/funk/internal/color"
)

var width = 11

func Get(name string) string {
	if width-len(name) < 0 {
		return name
	}
	return strings.Repeat(" ", width-len(name)) + name
}

var (
	ColorKind       = color.Green.Str(Get("kind"))
	ColorMsg        = color.Green.Str(Get("msg"))
	ColorService    = color.Green.Str(Get("service"))
	ColorOperation  = color.Green.Str(Get("operation"))
	ColorId         = color.Green.Str(Get("id"))
	ColorDetail     = color.Green.Str(Get("detail"))
	ColorTags       = color.Green.Str(Get("tags"))
	ColorErrMsg     = color.Red.Str(Get("err_msg"))
	ColorErrDetail  = color.Red.Str(Get("err_detail"))
	ColorCaller     = color.Green.Str(Get("caller"))
	ColorCode       = color.Green.Str(Get("code"))
	ColorMessage    = color.Green.Str(Get("message"))
	ColorBiz        = color.Green.Str(Get("biz_code"))
	ColorStatusCode = color.Green.Str(Get("status_code"))
	ColorName       = color.Green.Str(Get("name"))
	ColorStack      = color.Green.Str(Get("stack"))
	ColorVersion    = color.Green.Str(Get("version"))
)
