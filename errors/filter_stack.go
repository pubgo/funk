package errors

import (
	"strings"
	"sync"

	"github.com/pubgo/funk/stack"
)

var skipStack sync.Map

func RegStackFilter() {
	var s = stack.Caller(1)
	pkgL := strings.Split(s.Pkg, ".")
	skipStack.Store(strings.Join(pkgL[:len(pkgL)-1], "."), nil)
}

func init() {
	//RegStackFilter()
}
