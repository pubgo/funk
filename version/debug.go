package version

import (
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/recovery"
)

var mainPath string
var commitID string
var buildTime string
var version = "v0.0.1-dev-99"
var project string
var instanceID = uuid.New().String()

func init() {
	defer recovery.Exit(func() {
		pretty.Println(
			project,
			version,
			commitID,
			buildTime,
		)
	})

	bi, ok := debug.ReadBuildInfo()
	assert.If(!ok, "failed to read debug build info")

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

	assert.If(project == "", "project is null")
	assert.If(version == "", "version is null")
	assert.If(commitID == "", "commitID is null")
	assert.If(buildTime == "", "buildTime is null")
}
