package logx

import (
	"github.com/kr/pretty"
	"github.com/pubgo/x/q"
)

func Pretty(a ...interface{}) {
	Info(pretty.Sprint(a...))
}

func ColorPretty(args ...interface{}) {
	Info(string(q.Sq(args...)))
}
