# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced error handling with stack traces
- Result types for functional error handling
- Feature flag system with environment variable integration
- CLI flag generation for urfave/cli
- Generic utility functions for collections and comparisons

### Changed
- Refactored error wrapping mechanisms
- Improved performance of result type operations
- Enhanced documentation and examples

### Fixed
- Race conditions in feature flag registry
- Memory leaks in error stack traces
- Incorrect error wrapping in nested operations

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