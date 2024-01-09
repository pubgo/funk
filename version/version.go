package version

import (
	"runtime/debug"
)

var mainPath string

// git rev-parse HEAD
var commitID string
var buildTime string

// git describe --tags --abbrev=0
var version = "v0.0.1-dev-99"
var project = "project"

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
