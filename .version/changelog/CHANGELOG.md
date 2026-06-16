# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.4] - 2026-06-17

### Added
- `component/pyroscope`: Grafana Pyroscope continuous profiling integration (#62)

### Changed
- `component/cloudevent`: migrate catdogs runtime, split proto packages, and update `protoc-gen-go-cloudevent2` (#61)

## [2.0.3] - 2026-06-15

### Added
- Release tooling: git-cliff notes, GoReleaser v2 config, and `docs/RELEASE.md`
- CI: dedicated test job running `go test ./... -race`
- Makefile: `changelog`, `changelog-latest`, `release-dry`, and `install-tools` targets

### Fixed
- `async`: Iterator data race under `-race`; first-error-wins error retention; `Await` drains channel before returning
- `recovery`: injectable exit/fatal hooks for CI-safe tests
- Test fixes for `assert`, `async`, and `connmux` under the full CI race suite

### Changed
- GoReleaser pinned to v2.12.0; release workflow publishes git-cliff generated notes

## [2.0.2] - 2026-06-15

### Added
- `log`: `WithLogger`, `FromCtx`, and context-aware global helpers; `FlatMap` aliases in `result`
- `result`: `FlatMap` / `FlatMapTo` aliases; expanded regression tests (coverage ~80%+)
- CI: lint workflow now runs on `v2` branch

### Changed
- `log`: context field maps are cloned to avoid mutating logger defaults; `Logger.WithFields` caller fields win
- `result`: `Fail(nil)` panics as a programming mistake; `Recovery` calls `recover` directly

### Fixed
- `log`: `Error.MarshalJSON` encodes OK as JSON `null`; slog/std nil fallbacks
- `result`: `Error.MarshalJSON`, `Future.Await` race, `Async(nil)` deadlock, typed-nil `ErrSetter` propagation
- `result`: removed unsafe reflection from `setError`; concurrent-safe `resultchecker` registration

## [2.0.0] - 2025-12-06

### Added
- Major rewrite with generics support
- New result module with functional programming concepts
- Enhanced error module with better metadata support
- Features module for configuration management
- Comprehensive test suite with 90%+ coverage

### Changed
- Breaking API changes to support generics
- Improved performance across all modules
- Updated dependencies to latest versions
- Restructured package organization

### Removed
- Deprecated legacy APIs
- External dependency on third-party error libraries
- Redundant utility functions

## [1.0.0] - 2023-01-15

### Added
- Initial release
- Basic error handling utilities
- Simple result types
- Core utility functions
- Fundamental feature flag system