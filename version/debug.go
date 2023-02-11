package version

import (
	"runtime/debug"

	semver "github.com/hashicorp/go-version"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/recovery"
	"github.com/rs/xid"
)

var mainPath string
var commitID string
var buildTime string
var version = "v0.0.1-dev-99"
var project string
var instanceID = xid.New().String()

func init() {
	defer recovery.Exit(func(evt *errors.Event) {
		pretty.Println(
			mainPath,
			project,
			version,
			commitID,
			buildTime,
		)
	})

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

func Check() {
	defer recovery.Exit(func(evt *errors.Event) {
		pretty.Println(
			mainPath,
			project,
			version,
			commitID,
			buildTime,
		)
	})

	assert.Must1(semver.NewVersion(version))
	assert.If(project == "", "project is null")
	assert.If(version == "", "version is null")
	assert.If(commitID == "", "commitID is null")
	assert.If(buildTime == "", "buildTime is null")
}
