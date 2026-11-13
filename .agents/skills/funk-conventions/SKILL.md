---
name: funk-conventions
description: Funk 仓库核心约定速查——result/errors/log/env 助手优先、文档同步规则、覆盖率门禁、分支与 PR 流程。在 funk 仓库写代码、review、重构或回答仓库约定问题前使用。
---

# Funk 仓库约定速查

本 skill 是路由器:权威内容在仓库根的 `AGENTS.md`,这里只保留高频检查点。
**遇到细节冲突时,以 AGENTS.md 和工作区当前代码为准。**

## 动手前

1. 在仓库根目录执行所有 Go 命令。
2. import 保留 `/v2` 后缀。
3. 改动前先读目标包的 README(如存在)。

## 助手优先(不要手写重复实现)

- 日志:`log.Err(err)`(自动带 `error_chain`/`error_tags`);可读消息用 `errors.FormatChain(err)`。
- env 访问:`env.Get/GetOr/GetBoolOr/GetIntOr/GetFloatOr`、`env.Expand`
  (`${VAR:default}` 替换)——不要直接用 `os.Getenv`/`os.ExpandEnv`。
- 大小写转换:只用 vendored 的 `internal/strcase`(冻结,勿扩展、勿新增外部依赖)。

## 改动后

1. 导出行为变更 → 同步最近的 README / `_docs.go`(中英文都有则都改)。
2. 基础包(errors/result/log/config/env)→ 至少跑该包 + 邻近依赖包测试。
3. `golangci-lint run` 0 issues;覆盖率门禁(.testcoverage.yml)不许回退。

## 提交与 PR

- 新任务新分支;commit 用 Conventional Commits(`fix(env): ...`);不 push 除非明确要求。
- 含 `${...}` 的文本写进 commit message 时用 `-F <file>`,避免 shell 展开。
- push 后 CI 的 tidy-diff 会检查 go.mod/go.sum —— `go get` 之后必须 `go mod tidy`。
