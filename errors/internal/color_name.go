package internal

import (
	"strings"

	"github.com/pubgo/funk/internal/color"
)

var width = 10

func Get(name string) string {
	return strings.Repeat(" ", width-len(name)) + name
}

var ColorKind = color.Green.P(Get("kind"))
var ColorMsg = color.Green.P(Get("msg"))
var ColorDetail = color.Green.P(Get("detail"))
var ColorTags = color.Green.P(Get("tags"))
var ColorErrMsg = color.Red.P(Get("err_msg"))
var ColorErrDetail = color.Red.P(Get("err_detail"))
var ColorCaller = color.Green.P(Get("caller"))
var ColorCode = color.Green.P(Get("code"))
var ColorName = color.Green.P(Get("name"))
var ColorReason = color.Green.P(Get("reason"))
var ColorStatus = color.Green.P(Get("status"))
var ColorStack = color.Green.P(Get("stack"))
var ColorEvent = color.Green.P(Get("event"))
