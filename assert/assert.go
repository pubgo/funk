package assert

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"
)

func Assert(b bool, format string, a ...interface{}) {
	if b {
		must(fmt.Errorf(format, a...))
	}
}

func MustEqual[T any](a, b T) {
	if !cmp.Equal(a, b) {
		pp.Println("a: ", a)
		pp.Println("b: ", b)
		must(fmt.Errorf("a,b not equal"))
	}
}

func If(b bool, format string, a ...interface{}) {
	if b {
		must(fmt.Errorf(format, a...))
	}
}

func T(b bool, format string, a ...interface{}) {
	if b {
		must(fmt.Errorf(format, a...))
	}
}

func Err(b bool, err error) {
	if b {
		must(err)
	}
}

func Fn(b bool, fn func() error) {
	if b {
		must(fn())
	}
}

func Lazy(lazy func() bool, err error) {
	if lazy() {
		must(err)
	}
}
