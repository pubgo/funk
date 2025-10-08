package buildinfo

import (
	"github.com/samber/lo"
	"runtime/debug"
	"strings"
)

func CommitID() string  { return commitID }
func MainPath() string  { return mainPath }
func Version() string   { return version }
func BuildTime() string { return buildTime }
func Project() string   { return project }
func Domain() string    { return domain }

var domain string
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
	version string
	project string
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	mainPath = bi.Main.Path
	if project == "" {
		project = lo.LastOrEmpty(strings.Split(mainPath, "/"))
	}

	if version == "" {
		version = bi.Main.Version
	}

	if version == "" {
		version = "v0.0.1-dev-99"
	}

	for i := range bi.Settings {
		setting := bi.Settings[i]
		switch setting.Key {
		case "vcs.revision":
			commitID = setting.Value
		case "vcs.time":
			buildTime = setting.Value
		}
	}
}
