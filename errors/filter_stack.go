package errors

import (
	"strings"
	"sync"

	"github.com/pubgo/funk/stack"
)

var filterStack sync.Map

func RegStackFilter() {
	var s = stack.Caller(1)
	pkgL := strings.Split(s.Pkg, ".")
	filterStack.Store(strings.Join(pkgL[:len(pkgL)-1], "."), nil)
}

func init() {
	RegStackFilter()
}
