package version

import (
	"runtime/debug"
	"strings"
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

var (
	modified  bool
	os        string
	arch      string
	buildTags []string
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	mainPath = bi.Main.Path

	for i := range bi.Settings {
		setting := bi.Settings[i]
		switch setting.Key {
		case "vcs.revision":
			commitID = setting.Value
		case "vcs.time":
			buildTime = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		case "GOOS":
			os = setting.Value
		case "GOARCH":
			arch = setting.Value
		case "-tags":
			if setting.Value != "" {
				buildTags = strings.Split(setting.Value, ",")
			}
		}
	}
}
