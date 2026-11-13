---
description: 对指定 Go 包跑语句覆盖率并报告低覆盖函数与缺口
---

对 $ARGUMENTS 指定的包执行覆盖率检查(未给包名时,用 `git diff --name-only` 推断本次改动涉及的包):

1. 运行 `go test -coverprofile=/tmp/cover-check.out <pkgs>`(默认加 `-race`)。
2. 运行 `go tool cover -func=/tmp/cover-check.out`,输出:
   - 每个包的总覆盖率;
   - 覆盖率低于 70% 的函数清单(按覆盖率升序);
   - 与 `.testcoverage.yml` 中对应包的 override 阈值对比,标注是否贴近门禁。
3. 对 0% 覆盖的导出函数,给出一句补测建议(测什么分支),不要直接写测试代码。

报告用中文,表格呈现,最后给出一句总评。
