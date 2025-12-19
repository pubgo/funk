package buildinfo

var (
	domain   string
	mainPath string
)

// git rev-parse HEAD
// git describe --always --abbrev=7 --dirty
var (
	// commitID, git commit id
	commitID string

	// buildTime, build time, rfc3339
	buildTime string
)

// git describe --tags --abbrev=0
// git tag --sort=committerdate | tail -n 1
var (
	// version, git tag
	version string

	// project, project name
	project string
)
