package buildinfo

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/mod/module"

	v "github.com/pubgo/funk/v2/buildinfo/version"
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

	if module.IsPseudoVersion(bi.Main.Version) {
		ver := bi.Main.Version
		if a, err := module.PseudoVersionTime(ver); err == nil {
			buildTime = a.Format(time.RFC3339)
		}

		if b, err := module.PseudoVersionRev(ver); err == nil {
			commitID = b
		}

		if c, err := module.PseudoVersionBase(ver); err == nil {
			version = c
		}
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

	_ = v.SetBuildTime(buildTime)
	_ = v.SetCommitID(commitID)
	_ = v.SetProject(project)
	_ = v.SetMainPath(mainPath)
	_ = v.SetVersion(version)
	_ = v.SetDomain(domain)
}
