package version

import (
	"runtime/debug"
)

var mainPath string

// git rev-parse HEAD
// git describe --always --abbrev=7 --dirty
var (
	commitID  string
	buildTime string
)

// git describe --tags --abbrev=0
// git tag --sort=committerdate | tail -n 1
var (
	version = "v0.0.1-dev-99"
	project = "project"
)

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
