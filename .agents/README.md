# .agents — Agent 资产目录

团队共享的 agent 相关资产,随仓库版本化。各类编码工具对本目录的识别与加载规则
不同,以所用工具自己的文档为准。

## 目录约定

| 子目录 | 内容 |
| ------ | ---- |
| `skills/<name>/SKILL.md` | 能力包:按需触发的操作指南 |
| `commands/<name>.md` | prompt 模板(部分工具加载为斜杠命令)|
| `agents/<name>.md` | subagent 定义(是否被工具自动加载取决于该工具)|

## 现有资产

- skill: [`skills/funk-conventions/SKILL.md`](skills/funk-conventions/SKILL.md) — 仓库核心约定速查(路由到 AGENTS.md)
- command: [`commands/coverage-check.md`](commands/coverage-check.md) — 覆盖率检查
- agent: [`agents/docs-sync-checker.md`](agents/docs-sync-checker.md) — 文档同步检查(只读)

## 约定

1. 新增资产走 PR(与代码相同);不要提交密钥类内容。
2. skill 的 `description` 决定触发时机,必须写清"做什么 + 何时用"。
3. 路径一律使用仓库相对路径,不要写本机绝对路径。
