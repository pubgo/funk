# CloudEvent 模块

基于 NATS JetStream 的 protobuf CloudEvent 任务运行时。通过 `.proto` 中的 service/method 扩展声明 job 与 subject，并由 `protoc-gen-go-cloudevent2` 生成注册与发布代码。

## 功能

- **JetStream 任务**：通过 YAML 配置 stream、consumer、重试与并发
- **Protobuf 驱动**：job 名与 subject/topic 来自 proto 扩展
- **代码生成**：`Register*CloudEvent`、`*Publisher` 及 subject/job 常量
- **拦截器**：订阅端与发布端中间件链
- **投递控制**：`Reject`、`Redelivery`、`ForceRetry`
- **类型安全结果**：`Publish` 返回 `result.Result[*PubAckInfo]`

## 安装

```bash
go get github.com/pubgo/funk/v2/component/cloudevent
go install github.com/pubgo/funk/v2/cmds/protoc-gen-go-cloudevent2@latest
```

在本仓库根目录：

```bash
make protobuf
```

## Protobuf 结构

CloudEvent 刻意拆成 **两个** proto 包：

| 包 | Go import | 用途 |
|----|-----------|------|
| `proto/cloudevent` | `github.com/pubgo/funk/v2/proto/cloudevent`（`cloudeventpb`） | 消息类型：`PushEventOptions`、`RegisterJobOptions` 等 |
| `proto/cloudeventoption` | `github.com/pubgo/funk/v2/proto/cloudeventoption`（`cloudeventoptionpb`） | 描述符扩展 `E_Job`、`E_Subject` |

消息类型与扩展分离，业务 proto 只需 import 扩展包即可注解 service/method，而不必在每个生成 API 里混入扩展定义。

### 注解示例

```protobuf
syntax = "proto3";

import "cloudeventoption/options.proto";

service GidInnerService {
  option (lava.cloudevent.job) = { name: "gid" };

  rpc ProxyExec(ProxyExecReq) returns (google.protobuf.Empty) {
    option (lava.cloudevent.subject) = { name: "gid.proxy.exec" };
  }
}
```

- `job.name` 对应 YAML 中 `consumers` 下的 job 键名。
- `subject.name` 对应该 consumer 下的 subject 配置项。

### 代码生成

在 `protoc` 参数中加入（参见根目录 `protobuf.yaml`）：

```text
--go-cloudevent2_out=paths=source_relative:.
```

每个带注解的 service 会生成：

- `<Service>CloudEventJobKey` — job 名常量
- `<Method>CloudEventSubjectKey` — subject/topic 常量
- `<Service>CloudEvent` — 含 `On<Method>` 字段的 handler 结构体
- `Register<Service>CloudEvent` — 注册非 nil 的 handler
- `<Service>Publisher` — `{ Client, Opt, Interceptors }` 及 `Push<Method>Event` 方法

## 配置

任务在 YAML 中声明（结构见 `config.yaml`）：

```yaml
jobs:
  streams:
    gid:
      storage: "file"
      subjects: ["gid.>"]
  consumers:
    gid:
      - consumer: "test:gid"
        stream: "gid"
        subjects: "gid.proxy.exec"
        job:
          timeout: "1m"
          max_retries: 10
```

加载配置并绑定 NATS JetStream 连接创建 `Client`。配置中的 subject 在运行时会加上前缀（`DefaultPrefix`，默认 `acj`）。

默认项（可按 job 覆盖）：

| 项 | 默认值 |
|----|--------|
| Handler 超时 | `15s` |
| 最大重试 | `3` |
| 重试间隔 | `1s` |
| Consumer `AckWait` | `5m` |

## 快速开始

### 注册 handler（生成代码）

```go
RegisterGidInnerServiceCloudEvent(jobCli, GidInnerServiceCloudEvent{
    OnProxyExec: func(ctx context.Context, req *gidpb.ProxyExecReq) error {
        evt := cloudevent.GetContext(ctx)
        return nil
    },
})
```

### 注册 handler（手动）

```go
cloudevent.RegisterJobHandler(jobCli, "gid", "gid.proxy.exec",
    func(ctx context.Context, req *gidpb.ProxyExecReq) error { return nil },
    cloudevent.ProtoRegisterOpts(registerOpts...),
    cloudevent.WithSubInterceptors(myInterceptor),
)
```

### 发布（生成 Publisher）

```go
pub := GidInnerServicePublisher{
    Client: jobCli,
    Opt:    cloudevent.ProtoPubOpts(defaultPushOpts...),
}
ack := pub.PushProxyExecEvent(ctx, req).Unwrap()
```

### 发布（直接调用）

```go
ack := cloudevent.Publish(jobCli, ctx, subjectKey, req, interceptors, opts...).Unwrap()
```

`Publish` 会合并 `PubOpt`、默认 `content_type` 为 `application/json`，未指定时自动填充 `sender`。

## 拦截器

**订阅端**（`SubInterceptor`）：包装 handler 执行，通过 `WithSubInterceptors` 或 `RegisterOpt` 传入。

**发布端**（`PubInterceptor`）：包装发布逻辑，通过 `Publisher.Interceptors` 或 `Publish` 参数传入。

均采用 `func(next ...) ...` 中间件形式。

## 投递错误

Handler 返回以下错误可控制 JetStream 确认行为：

| 函数 | 行为 |
|------|------|
| `Reject(errs...)` | Ack 并丢弃，不再重试 |
| `Redelivery(delay, errs...)` | Nak 并延迟重投 |
| `ForceRetry(errs...)` | 立即 Nak，强制重投 |

普通 error 按配置的 retry/backoff 重试，直至达到 `max_retries`。

## Context

`GetContext(ctx)` 返回投递元数据：subject、stream、consumer、投递次数、时间戳、HTTP 风格 header 及 `JobEventConfig`。旧 API `GetEventContext` 已更名为 `GetContext`。

## 从旧版 funk cloudevent 迁移

| 旧版 | 新版 |
|------|------|
| `proto/cloudevent/options.proto` 内联扩展 | `proto/cloudeventoption/options.proto`（`cloudeventoptionpb`） |
| `GetEventContext` | `GetContext` |
| `PushEvent` / `PushRpcEvent` | 生成的 `*Publisher.Push*Event` 或 `Publish` |
| `Publish` 返回 `(ack, error)` | `result.Result[*PubAckInfo]` |
| 无 `RegisterOpt` | `RegisterJobHandler(..., opts ...RegisterOpt)` |
| 单一 proto 包 | `cloudeventpb` + `cloudeventoptionpb` |

下游服务迁移清单：

1. proto import 改为 `cloudeventoption/options.proto`
2. 使用 `protoc-gen-go-cloudevent2` 重新生成
3. 用 `Register*CloudEvent` 或带 `RegisterOpt` 的手动注册替换旧注册方式
4. 用 `*Publisher` 或 `Publish` + `result` 替换旧发布 helper
5. `GetEventContext` 改名为 `GetContext`
6. YAML job/subject 与 proto 注解保持一致

## 参考

- [CloudEvents Go SDK](https://github.com/cloudevents/sdk-go)
- [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
