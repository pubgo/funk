package version

import (
	"runtime"

	vv "github.com/hashicorp/go-version"
	"github.com/pubgo/funk/assert"
)

func init() {
	assert.Exit1(vv.NewVersion(runtime.Version()))
}
