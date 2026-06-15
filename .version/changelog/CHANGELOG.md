# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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