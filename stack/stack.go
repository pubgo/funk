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

func CallerWithDepth(skip int) string {
	var pcs [1]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return ""
	}

	return stack(pcs[0] - 1).String()
}

func Caller(skip int) *Func {
	var pcs [1]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return nil
	}

	return stack(pcs[0] - 1)
}

func Callers(depth int) []*Func {
	var pcs = make([]uintptr, depth)
	n := runtime.Callers(2, pcs[:])
	if n == 0 {
		return nil
	}

	var stacks = make([]*Func, 0, depth)
	for _, p := range pcs[:n] {
		stacks = append(stacks, stack(p-1))
	}
	return stacks
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
