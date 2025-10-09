package version

import (
	_ "embed"
)

//go:embed .version
var version string

// ReleaseVersion v2.0.0
func ReleaseVersion() string { return version }

// ReleaseDate 2025-10-10T12:01:46Z
func ReleaseDate() int64 { return 1760097706 }
