package version

import (
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/pubgo/funk/assert"
)

func init() {
	assert.Exit1(version.NewVersion(runtime.Version()))
}
