package errinter

import (
	"strings"

	"github.com/pubgo/funk/v2/internal/color"
)

var width = 11

func getColorErrField(name string) string {
	if width-len(name) < 0 {
		return name
	}
	return strings.Repeat(" ", width-len(name)) + name
}

var (
	ColorKind       = color.Green.Str(getColorErrField("kind"))
	ColorMsg        = color.Green.Str(getColorErrField("msg"))
	ColorService    = color.Green.Str(getColorErrField("service"))
	ColorOperation  = color.Green.Str(getColorErrField("operation"))
	ColorId         = color.Green.Str(getColorErrField("id"))
	ColorDetail     = color.Green.Str(getColorErrField("detail"))
	ColorTags       = color.Green.Str(getColorErrField("tags"))
	ColorErrMsg     = color.Red.Str(getColorErrField("err_msg"))
	ColorErrDetail  = color.Red.Str(getColorErrField("err_detail"))
	ColorCaller     = color.Green.Str(getColorErrField("caller"))
	ColorCode       = color.Green.Str(getColorErrField("code"))
	ColorMessage    = color.Green.Str(getColorErrField("message"))
	ColorBiz        = color.Green.Str(getColorErrField("biz_code"))
	ColorStatusCode = color.Green.Str(getColorErrField("status_code"))
	ColorName       = color.Green.Str(getColorErrField("name"))
	ColorStack      = color.Green.Str(getColorErrField("stack"))
	ColorVersion    = color.Green.Str(getColorErrField("version"))
)
