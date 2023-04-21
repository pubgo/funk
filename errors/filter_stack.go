package errors

import (
	"strings"
	"sync"

	"github.com/pubgo/funk/stack"
)

var skipStackMap sync.Map

func RegStackPkgFilter() {
	var s = stack.Caller(1)
	pkgL := strings.Split(s.Pkg, ".")
	skipStackMap.Store(strings.Join(pkgL[:len(pkgL)-1], "."), nil)
}
