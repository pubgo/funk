package version

import "testing"

func TestVersionUsesEnvOverride(t *testing.T) {
	t.Setenv(envVersion, "v9.9.9-test")
	version = "v0.0.1"

	if got := Version(); got != "v9.9.9-test" {
		t.Fatalf("Version() = %q, want %q", got, "v9.9.9-test")
	}
}

func TestVersionFallsBackWithoutEnv(t *testing.T) {
	version = "v1.2.3"

	if got := Version(); got != "v1.2.3" {
		t.Fatalf("Version() = %q, want %q", got, "v1.2.3")
	}
}

func TestCommitIDUsesEnvOverride(t *testing.T) {
	t.Setenv(envCommitID, "commit-from-env")
	commitID = "commit-from-build"

	if got := CommitID(); got != "commit-from-env" {
		t.Fatalf("CommitID() = %q, want %q", got, "commit-from-env")
	}
}

func TestBuildTimeUsesEnvOverride(t *testing.T) {
	t.Setenv(envBuildTime, "2026-05-15T12:00:00Z")
	buildTime = "2026-05-15T11:00:00Z"

	if got := BuildTime(); got != "2026-05-15T12:00:00Z" {
		t.Fatalf("BuildTime() = %q, want %q", got, "2026-05-15T12:00:00Z")
	}
}

func TestProjectUsesEnvOverride(t *testing.T) {
	t.Setenv(envProject, "portal")
	project = "project"

	if got := Project(); got != "portal" {
		t.Fatalf("Project() = %q, want %q", got, "portal")
	}
}
