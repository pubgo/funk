package stack

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var (
	cache  sync.Map
	mu     sync.Mutex
	goRoot string
)

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
	ff := f.File
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
	skip := 0
	if len(skips) > 0 {
		skip = skips[0]
	}

	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return nil
	}

	stacks := make([]*Frame, 0, depth)
	for _, p := range pcs[:n] {
		stacks = append(stacks, stack(p-1))
	}
	return stacks
}

func CallerWithType(typ reflect.Type) *Frame {
	if typ == nil {
		return nil
	}

	return &Frame{Pkg: typ.PkgPath(), Name: typ.Name(), File: typ.PkgPath()}
}

func CallerWithFunc(fn interface{}) *Frame {
	if fn == nil {
		panic("[fn] param is nil")
	}

	var vfn reflect.Value
	if v, ok := fn.(reflect.Value); ok {
		vfn = v
	} else {
		vfn = reflect.ValueOf(fn)
	}

	if !vfn.IsValid() || vfn.Kind() != reflect.Func || vfn.IsNil() {
		panic("[fn] is not func type or type is nil")
	}

	return stack(vfn.Pointer())
}

func GetStack(skip int) uintptr {
	var pcs [1]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	if n == 0 {
		return 0
	}

	return pcs[0] - 1
}

func Stack(p uintptr) *Frame {
	return stack(p)
}

func stack(p uintptr) *Frame {
	if p == 0 {
		return nil
	}

	v, ok := cache.Load(p)
	if ok {
		return v.(*Frame)
	}

	mu.Lock()
	defer mu.Unlock()

	v, ok = cache.Load(p)
	if ok {
		return v.(*Frame)
	}

	defer func() { cache.Store(p, v) }()

	ff := runtime.FuncForPC(p)
	if ff == nil {
		return nil
	}

	file, line := ff.FileLine(p)
	ma := strings.Split(ff.Name(), ".")
	v = &Frame{
		File: file,
		Line: line,
		Name: ma[len(ma)-1],
		Pkg:  strings.Join(ma[:len(ma)-1], "."),
	}
	return v.(*Frame)
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
