package version

import (
	"runtime/debug"

	"github.com/rs/xid"
)

var mainPath string
var commitID string
var buildTime string
var version = "v0.0.1-dev-99"
var project string
var instanceID = xid.New().String()

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

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
}
