# .agents — Agent 资产目录

团队共享的 agent 相关资产,随仓库版本化。ZCode 会在会话启动时自动扫描本目录的
部分子目录(规则见下),其余工具(Claude Code、Codex 等)也可以读取同一份资产。

## 目录约定与加载规则

| 子目录 | 内容 | ZCode 是否自动加载 |
| ------ | ---- | ------------------ |
| `skills/<name>/SKILL.md` | 能力包:按需触发的操作指南 | ✅ 自动发现(触发由 description 决定) |
| `commands/<name>.md` | prompt 模板,加载为 `/name` 斜杠命令 | ✅ 自动注册为命令 |
| `prompts/` | 非 slash 命令形态的可复用提示词文本 | ❌ 仅供人/其他工具取用 |

> ZCode 同时扫描 `.zcode/` 与 `.agents/` 两个位置,同名资源 `.zcode` 优先。
> skill 触发、命令重名等问题的排查见 `zcode-guide:diagnosing-skills` / `diagnosing-commands`。

## Subagent(子代理)单独放置

ZCode 的 subagent 定义**只**从 `<repo>/.zcode/agents/*.md` 加载(不在 `.agents/` 扫描范围),
因此本仓库的子代理统一放在 [`../.zcode/agents/`](../.zcode/agents/)。
frontmatter 携带 `name` / `description` / `tools`,正文为子代理的系统指令。

## 现有资产

- skill: [`skills/funk-conventions/SKILL.md`](skills/funk-conventions/SKILL.md) — 仓库核心约定速查(路由到 AGENTS.md)
- command: [`commands/coverage-check.md`](commands/coverage-check.md) — `/coverage-check <pkg>` 覆盖率报告
- subagent: [`../.zcode/agents/docs-sync-checker.md`](../.zcode/agents/docs-sync-checker.md) — 文档同步检查

## 约定

1. 新增资产走 PR(与代码相同);不要提交密钥类内容。
2. skill 的 `description` 决定触发时机,必须写清"做什么 + 何时用"。
3. 命令重名时:用户级(`~/.zcode`、`~/.agents`)优先于仓库级,先查重再新增。
