package debug

import (
	"runtime/debug"

	"github.com/kr/pretty"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/recovery"
)

var commitID string
var buildTime string
var mainPath string

func init() {
	defer recovery.Exit(func() {
		pretty.Log(commitID, buildTime)
	})

	bi, ok := debug.ReadBuildInfo()
	assert.If(!ok, "failed to read build info")

	mainPath = bi.Main.Path

	for i := range bi.Settings {
		setting := bi.Settings[i]
		if setting.Key == "vcs.revision" {
			commitID = setting.Value
		}

		if setting.Key == "vcs.time" {
			buildTime = setting.Value
		}
	}

	assert.If(commitID == "", "commitID is null")
	assert.If(buildTime == "", "buildTime is null")
}
