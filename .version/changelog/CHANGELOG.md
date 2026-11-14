# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.0] - 2026-11-14

### Added
- `env`: `Expand` supports `${VAR:default}` / `${VAR:-default}` (default applies when the variable is unset or empty), quote-stripped defaults and `$$` escaping; substitution is implemented in-repo
- `env`: `GetBoolOr`, `GetIntOr`, `GetFloatOr` typed getters with explicit defaults
- `env`: `MustLoadFiles`, the panicking variant of `LoadFiles`
- `log`: `SetDisabledLoggers` glob-based logger suppression (dot-separated names, bare patterns include their subtree) with `IsDisabled` / `IsDisabledEvent`
- CI: go-test-coverage gate (`.testcoverage.yml`) and docs publishing via pubgo/markview (GitHub Pages)

### Changed
- `env`/`config`: environment substitution no longer uses `github.com/a8m/envsubst`; keys resolve raw-name-first with normalized fallbacks, and `Delete` unsets every stored spelling of a key
- `env`: `Get` and the typed getters trim values and follow the current `slog.Default()`
- vendored `github.com/ettle/strcase` as `internal/strcase` (frozen, MIT); dropped `github.com/iancoleman/strcase`
- `cmds/protoc-gen-go-errors2`: naming now uses the vendored strcase — registered identity strings change for names with digits or acronym boundaries (e.g. `sha_256_mismatch` -> `sha256_mismatch`)
- `component/resty`: dropped deprecated `net.Dialer.DualStack`
- `.agents/` directory for team-shared agent assets (skills, commands)

### Fixed
- `env`: `${VAR:default}` applies defaults again — the previous envsubst behavior treated the single-colon form as a bash substring offset, so `${PORT:8080}` resolved to empty
- `env`: keys containing digits resolve regardless of case (`${sha256sum}` finds `SHA256_SUM` stored via `env.Set`)
- `log`: hooks and checkers see the component name through the context field map
- `component/resty`: staticcheck SA1019 lint failure

### Removed
- `github.com/a8m/envsubst` and `github.com/iancoleman/strcase` dependencies

## [2.0.5] - 2026-06-25

### Added
- `errors`: `GetTags`, `CollectTags`, `CollectUserTags`, `FormatChain`, `FullMessage`, `MarshalError`, `Walk`, and `RootCause`
- `errors/errcode`: `RegisterErrCode`, `MustRegisterErrCode`, `LookupErrCode`, and package README
- `log`: `Err()` now emits `error_chain` and `error_tags` automatically
- `result`: `Message()` and `Tags()` helpers on `Error`

### Changed
- `errors`: `ErrJsonify` now expands the full unwrap chain; `WrapStack` no longer prints to stderr; `JsonPrint` returns nil instead of panicking on marshal failure
- `errors`: tag maps are cloned on create/wrap to avoid accidental mutation
- `result`: `String()` uses the formatted error chain for clearer output
- `errors/errcode`: `NewCodeErrWithMsg` preserves caller casing instead of using `strings.ToTitle`

### Fixed
- `errors/errcode`: `ErrCode.As` correctly extracts `*ErrCode` and `*errorpb.ErrCode`
- `errors/errcode`: duplicate registration panic message now reports the conflicting name
- `errors`: `Unwrap` falls back to the standard library unwrap behavior

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