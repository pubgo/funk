package stack

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var cache = make(map[uintptr]*Frame)
var mu sync.Mutex
var goRoot string

func init() {
	tt := Trace()[0]
	goRoot = tt.File[:pkgIndex(tt.File, tt.Pkg)]
}

type Frame struct {
	Name string `json:"name"`
	Pkg  string `json:"pkg"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (f *Frame) Short() string {
	var ff = f.File
	ff1 := ff[:strings.LastIndex(ff, "/")]
	ff1 = ff1[:strings.LastIndex(ff1, "/")]
	return fmt.Sprintf("%s:%d %s", strings.TrimPrefix(strings.TrimPrefix(ff, ff1), "/"), f.Line, f.Name)
}

func (f *Frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Name)
}

func (f *Frame) IsRuntime() bool {
	return strings.Contains(f.File, goRoot)
}

func GetGORoot() string { return goRoot }

func Caller(skip int) *Frame {
	var pcs [1]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return nil
	}

	return stack(pcs[0] - 1)
}

func Callers(depth int, skips ...int) []*Frame {
	var skip = 0
	if len(skips) > 0 {
		skip = skips[0]
	}

	var pcs = make([]uintptr, depth)
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return nil
	}

	var stacks = make([]*Frame, 0, depth)
	for _, p := range pcs[:n] {
		stacks = append(stacks, stack(p-1))
	}
	return stacks
}

func CallerWithFunc(fn interface{}) *Frame {
	if fn == nil {
		return nil
	}

	var fn1 = reflect.ValueOf(fn)
	if !fn1.IsValid() || fn1.Kind() != reflect.Func || fn1.IsNil() {
		panic("[fn] is not func type or type is nil")
	}

	return stack(fn1.Pointer())
}

func stack(p uintptr) *Frame {
	var v, ok = cache[p]
	if ok {
		return v
	}

	mu.Lock()
	defer mu.Unlock()

	v, ok = cache[p]
	if ok {
		return v
	}

	defer func() {
		cache[p] = v
	}()

	var ff = runtime.FuncForPC(p)
	if ff == nil {
		return nil
	}

	var file, line = ff.FileLine(p)
	ma := strings.Split(ff.Name(), ".")
	v = &Frame{File: file, Line: line, Name: ma[len(ma)-1], Pkg: strings.Join(ma[:len(ma)-1], ".")}
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
