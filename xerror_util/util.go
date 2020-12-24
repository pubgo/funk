package xerror_util

import (
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/pubgo/xerror/xerror_envs"
)

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func CallerWithDepth(callDepths ...int) string {
	if !xerror_envs.IsCallerVal() {
		return ""
	}

	var cd = xerror_envs.CallDepthVal()
	if len(callDepths) > 0 {
		cd = callDepths[0]
	}

	var pcs = make([]uintptr, 1)
	if runtime.Callers(cd, pcs[:]) == 0 {
		return ""
	}

	f := frame(pcs[0])
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown type"
	}

	file, line := fn.FileLine(f.pc())
	return file + ":" + strconv.Itoa(line)
}

func CallerWithFunc(fn interface{}) string {
	if fn == nil {
		log.Fatalln("params is nil")
	}

	var _fn = reflect.ValueOf(fn)
	if !_fn.IsValid() || _fn.IsNil() || _fn.Kind() != reflect.Func {
		log.Fatalln("not func type or type is nil")
	}

	var _e = runtime.FuncForPC(_fn.Pointer())
	var file, line = _e.FileLine(_fn.Pointer())

	var buf = &strings.Builder{}
	defer buf.Reset()

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")

	ma := strings.Split(_e.Name(), ".")
	buf.WriteString(ma[len(ma)-1])
	return buf.String()
}
