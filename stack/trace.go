package stack

import "runtime"

func Trace() []*Func {
	var pcs [512]uintptr
	n := runtime.Callers(1, pcs[:])
	cs := make([]*Func, 0, n)

	for _, p := range pcs[:n] {
		cs = append(cs, stack(p))
	}

	return cs
}
