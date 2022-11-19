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
	// As of Go 1.9 we need room for up to three PC entries.
	//
	// 0. An entry for the stack frame prior to the target to check for
	//    special handling needed if that prior entry is runtime.sigpanic.
	// 1. A possible second entry to hold metadata about skipped inlined
	//    functions. If inline functions were not skipped the target frame
	//    PC will be here.
	// 2. A third entry for the target frame PC when the second entry
	//    is used for skipped inline functions.
	var pcs [3]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	if n == 0 {
		return ""
	}

	pcs1 := pcs[:n]
	pc := pcs1[len(pcs1)-1]

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown type"
	}

	file, line := fn.FileLine(pc)
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
