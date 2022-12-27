package stack

import (
	"runtime"
)

func Trace() []*Frame {
	var pcs [512]uintptr
	n := runtime.Callers(0, pcs[:])
	cs := make([]*Frame, 0, n)

	for _, p := range pcs[:n] {
		cs = append(cs, stack(p))
	}

	return cs
}
