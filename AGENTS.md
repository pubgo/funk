# AGENTS.md

This repository is a Go utility library module: `github.com/pubgo/funk/v2`.

## What this repo contains

- Many small, reusable Go packages under the repository root.
- Core modules are documented in their own READMEs:
  - [root overview](./README.md)
  - [design](./docs/DESIGN.md)
  - [errors](./errors/README.md)
  - [errors/errcode](./errors/errcode/README.md)
  - [result](./result/README.md)
  - [config](./config/README.md)
  - [features](./features/README.md)
  - [log](./log/README.md)
  - [env](./env/README.md)
  - [stack](./stack/README.md)
  - [connmux](./connmux/README.md)
  - [cloudevent](./component/cloudevent/README.md)
  - [pyroscope](./component/pyroscope/README.md)

## Working rules for AI coding agents

1. Always run Go commands from the repository root: `/Users/barry/git/funk`.
2. Preserve the module import path suffix `/v2` when adding or editing imports.
3. Prefer the smallest possible change; do not reformat unrelated packages.
4. Before editing a package, read that package README first if it exists.
5. For package-local changes, prefer targeted tests before running wider suites.
6. After changing exported behavior, update the nearest README or package docs.
7. Treat logger/context field maps as immutable data; avoid mutating reused state in place.
8. Prefer existing helpers from `errors`, `result`, `log`, `config`, and `features` over introducing duplicate patterns.
   - For logs: use `log.Err(err)` instead of manually formatting chains; it emits `error_chain` and `error_tags`.
   - For readable messages: use `errors.FormatChain(err)` rather than concatenating `err.Error()` with wrap context.
   - For metadata: use `errors.CollectUserTags(err)` instead of reading wrap-layer `msg` tags directly.
9. Keep generated/proto-related changes scoped; do not edit generated output casually.
10. If a change touches shared APIs, search for cross-package usages before refactoring.

## Common commands

Run from repo root:

- Full test + race + coverage profile via Makefile: `make test`
- Coverage HTML: `make test_html`
- Benchmarks (root package target from Makefile): `make test_bench`
- All-package vet: `make vet`
- Lint: `make lint`
- Format/refactor: `make refactor`
- Release notes preview: `make changelog` / `make changelog-latest` (requires [git-cliff](https://git-cliff.org))
- Release dry-run: `make release-dry` (requires GoReleaser)
- Release process: [docs/RELEASE.md](./docs/RELEASE.md)
- All-package Go tests directly: `go test ./...`
- Package-scoped tests: `go test ./log`, `go test ./errors`, etc.

## Protobuf / code generation

Proto generation is configured in [protobuf.yaml](./protobuf.yaml).

Relevant commands:

- `make protobuf`
- `make protolint`

Notes:

- Generated output goes to `./proto`
- Vendored proto dependencies live in `./proto-vendor`
- Custom plugins live under `./cmds/`

## Repository map

- `errors/`: enriched error model, wrapping, tags, stack-aware debugging
- `result/`: generic `Result[T]` / `Error` flow helpers and async helpers
- `log/`: zerolog-based structured logging facade, context-aware fields, slog/std adapters
- `config/`: YAML config loading, env substitution, merge/patch support, expression helpers
- `features/`: typed feature flags, env + CLI integration
- `env/`: environment variable helpers and `.env` support
- `stack/`: caller and stack-trace helpers
- `async/`: future/promise/iterator helpers
- `component/`: integrations/adapters (bbolt, gorm, nats, cloudevent, etc.)
- `cmds/`: codegen and upgrade tools
- `docs/`: high-level design docs

## Testing guidance

- Start with the most local package affected.
- Prefer focused regression tests for bug fixes.
- If changing behavior in foundational packages like `errors`, `result`, `log`, or `config`, run at least that package and nearby dependents.
- For performance-sensitive changes, add or run benchmarks near the edited package when practical.

## Known repo conventions

- Package docs often live in `_doc.go`, `_docs.go`, `README.md`, and `README.zh.md`.
- Many packages are intentionally lightweight facades over existing libraries; preserve that style.
- The logging package intentionally depends on `zerolog`; do not abstract it away unless explicitly requested.
- Configuration and feature helpers favor typed APIs and declarative usage over ad-hoc maps.

## When in doubt

- Prefer linking to existing docs instead of copying them into new files.
- Follow the style and naming already present in the touched package.
- If a Make target and real package layout disagree, trust the current package structure and verify with `go test` / `go vet`.
