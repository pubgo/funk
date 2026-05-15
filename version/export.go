package version

import "os"

const (
	envProject   = "FUNK_PROJECT"
	envVersion   = "FUNK_VERSION"
	envCommitID  = "FUNK_COMMIT_ID"
	envBuildTime = "FUNK_BUILD_TIME"
)

func CommitID() string {
	return lookupEnv(envCommitID, commitID)
}

func MainPath() string {
	return mainPath
}

func Version() string {
	return lookupEnv(envVersion, version)
}

func BuildTime() string {
	return lookupEnv(envBuildTime, buildTime)
}

func Project() string {
	return lookupEnv(envProject, project)
}

func lookupEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
