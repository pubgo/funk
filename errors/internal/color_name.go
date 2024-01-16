package internal

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

var ColorKind = color.Green.P(Get("kind"))
var ColorMsg = color.Green.P(Get("msg"))
var ColorService = color.Green.P(Get("service"))
var ColorOperation = color.Green.P(Get("operation"))
var ColorId = color.Green.P(Get("id"))
var ColorDetail = color.Green.P(Get("detail"))
var ColorTags = color.Green.P(Get("tags"))
var ColorErrMsg = color.Red.P(Get("err_msg"))
var ColorErrDetail = color.Red.P(Get("err_detail"))
var ColorCaller = color.Green.P(Get("caller"))
var ColorCode = color.Green.P(Get("code"))
var ColorReason = color.Green.P(Get("reason"))
var ColorBiz = color.Green.P(Get("biz_code"))
var ColorStatusCode = color.Green.P(Get("status_code"))
var ColorName = color.Green.P(Get("name"))
var ColorStack = color.Green.P(Get("stack"))
var ColorVersion = color.Green.P(Get("version"))
