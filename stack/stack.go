package stack

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/cornelk/hashmap"
)

var cache = hashmap.New[uintptr, *Func]()
var goRoot string

func init() {
	tt := Trace()[0]
	goRoot = tt.File[:pkgIndex(tt.File, tt.Pkg)]
}

type Func struct {
	Name string `json:"name"`
	Pkg  string `json:"pkg"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (f *Func) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Name)
}

func (f *Func) TrimRuntime() *Func {
	var f1 = *f
	f1.File = strings.TrimPrefix(f1.File, goRoot)
	return &f1
}

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) }

func CallerWithDepth(cd int) string {
	var pcs = make([]uintptr, 1)
	if runtime.Callers(cd+2, pcs[:]) == 0 {
		return ""
	}

	f := frame(pcs[0])
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown type"
	}

	file, line := fn.FileLine(f.pc())
	return fmt.Sprintf("%s:%d", file, line)
}

func CallerWithFunc(fn interface{}) string {
	if fn == nil {
		return ""
	}

	var _fn = reflect.ValueOf(fn)
	if !_fn.IsValid() || _fn.Kind() != reflect.Func || _fn.IsNil() {
		panic("[fn] is not func type or type is nil")
	}

	var _e = runtime.FuncForPC(_fn.Pointer())
	var file, line = _e.FileLine(_fn.Pointer())

	ma := strings.Split(_e.Name(), ".")
	return fmt.Sprintf("%s:%d %s", file, line, ma[len(ma)-1])
}

func stack(p uintptr) *Func {
	var v, ok = cache.Get(p)
	if ok {
		return v
	}

	var ff = runtime.FuncForPC(p)
	if ff == nil {
		return nil
	}

	var file, line = ff.FileLine(p)
	ma := strings.Split(ff.Name(), ".")
	v = &Func{
		File: file,
		Line: line,
		Name: ma[len(ma)-1],
		Pkg:  strings.Join(ma[:len(ma)-1], "."),
	}
	cache.Set(p, v)
	return v
}

func pkgIndex(file, funcName string) int {
	const sep = "/"
	i := len(file)
	for n := strings.Count(funcName, sep) + 2; n > 0; n-- {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			i = -len(sep)
			break
		}
	}

	return i + len(sep)
}
