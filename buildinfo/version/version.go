package version

import "github.com/pubgo/funk/v2"

var (
	mainPath  string
	domain    string
	commitID  string
	buildTime string
	version   string
	project   string
	release   string
)

func CommitID() string       { return commitID }
func MainPath() string       { return mainPath }
func Version() string        { return version }
func ReleaseVersion() string { return release }
func BuildTime() string      { return buildTime }
func Project() string        { return project }
func Domain() string         { return domain }

func SetCommitID(id string) funk.Void {
	if id != "" {
		commitID = id
	}
	return funk.Void{}
}

func SetMainPath(path string) funk.Void {
	if path != "" {
		mainPath = path
	}
	return funk.Void{}
}

func SetVersion(v string) funk.Void {
	if v != "" {
		version = v
	}
	return funk.Void{}
}

func SetReleaseVersion(r string) funk.Void {
	if r != "" {
		release = r
	}
	return funk.Void{}
}

func SetBuildTime(t string) funk.Void {
	if t != "" {
		buildTime = t
	}
	return funk.Void{}
}

func SetProject(p string) funk.Void {
	if p != "" {
		project = p
	}
	return funk.Void{}
}

func SetDomain(d string) funk.Void {
	if d != "" {
		domain = d
	}
	return funk.Void{}
}
