# Release Guide

Funk is primarily a **Go module** (`github.com/pubgo/funk/v2`) and secondarily ships **protoc plugin binaries** via GitHub Releases.

This guide follows the same principles as [zigo](https://github.com/pubgo/zigo): **tags are the source of truth**, changelogs are generated from git history, and CI performs the release.

## Version model

| Artifact | Version source | Example |
|----------|----------------|---------|
| Go module | Git tag on `v2` | `v2.0.2` → `go get github.com/pubgo/funk/v2@v2.0.2` |
| Protoc plugins | Same tag + GoReleaser | `protoc-gen-go-errors-v2.0.2-darwin-arm64` |
| `.version/VERSION` | Release prep commit (human-facing) | `v2.0.2` |

Unlike zigo's per-tool tags (`mqttcli/v*`), funk keeps a **single module tag** so Go proxy consumers stay aligned with binary releases.

## Channels

| Tag suffix | GitHub Release | Example |
|------------|----------------|---------|
| (none) | Stable | `v2.0.2` |
| `-alpha` | Prerelease | `v2.0.3-alpha.1` |
| `-beta` | Prerelease | `v2.0.3-beta.1` |

GoReleaser sets `prerelease: auto` based on the semver suffix.

## Cut a release

### 1. Preflight on `v2`

```bash
make test
make vet
make lint
```

### 2. Preview changelog

```bash
make changelog          # unreleased commits
make changelog-latest   # since last tag
```

### 3. Prepare version files

Update:

- `.version/VERSION`
- `.version/changelog/CHANGELOG.md` (move `[Unreleased]` entries into the new version section)

```bash
git add .version/
git commit -m "chore(release): prepare v2.0.3"
```

### 4. Tag and push

```bash
git tag -a v2.0.3 -m "v2.0.3"
git push origin v2
git push origin v2.0.3
```

Pushing a `v*.*.*` tag triggers `.github/workflows/release.yml` (GoReleaser + git-cliff notes).

### 5. Verify

- GitHub Actions: **Release** workflow succeeds
- Release page contains plugin binaries + checksums
- Module resolves: `go get github.com/pubgo/funk/v2@v2.0.3`

## Local dry-run

```bash
make release-dry
```

Builds snapshot binaries under `dist/` without publishing.

## What we adopted from zigo

- **git-cliff** for scoped, conventional-commit release notes
- **Tag-driven CI** with prerelease channel detection
- **Makefile for dev**, workflows for release
- **Documented release checklist** (this file)

## What we intentionally differ on

- **Single `v2.x.y` tag** instead of per-binary tags (Go module requirement)
- **GoReleaser** instead of custom cross-compile scripts (Go toolchain)
- **Keep a Changelog** file in `.version/changelog/` for long-form history
