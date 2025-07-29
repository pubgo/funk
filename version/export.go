package version

func CommitID() string {
	return commitID
}

func MainPath() string {
	return mainPath
}

func Version() string {
	return version
}

func BuildTime() string {
	return buildTime
}

func Project() string {
	return project
}

func SetVersion(v string) { version = v }
func SetProject(p string) { project = p }
