# connmux

`connmux` 是一个基于连接“前若干字节”进行协议分流的连接复用器：你可以在同一个端口上同时跑 gRPC、HTTP/1、HTTP/2 或自定义 TCP 协议。

> 设计理念与 `github.com/soheilhy/cmux` 类似，但这里是面向本仓库的独立实现；同时提供 `github.com/pubgo/funk/v2/cmux` 兼容包装层，便于平滑迁移。

## 快速开始

```go
root, _ := net.Listen("tcp", ":8080")

m := connmux.New(root,
	connmux.WithReadTimeout(2*time.Second),
	connmux.WithMaxSniffBytes(1<<20),
)

// 按注册顺序匹配，越早注册优先级越高
grpcL := m.Match(connmux.HTTP2HeaderField("content-type", "application/grpc"))
httpL := m.Match(connmux.HTTP1Fast())
rawL := m.Match(connmux.Any())

go grpcServer.Serve(grpcL)
go httpServer.Serve(httpL)
go serveRaw(rawL)

_ = m.Serve()
```

## API 概览

- `New(root net.Listener, opts ...Option) *Mux`
- `(*Mux).Match(matchers ...Matcher) net.Listener`
- `(*Mux).MatchWithWriters(writers ...MatchWriter) net.Listener`
- `(*Mux).Serve() error`
- `(*Mux).Close() error`

### Options

- `WithReadTimeout(d time.Duration)`：sniff 阶段每次 Read 的超时；防 slowloris
- `WithMaxSniffBytes(n int)`：sniff 阶段最大缓存字节数（每连接）
- `WithConnBacklog(n int)`：每个子 listener 的队列长度（缓冲 Accept）
- `WithErrorHandler(func(error) bool)`：错误处理；返回 `true` 表示继续 Serve

## 内置匹配器

- `Any()`：兜底匹配
- `Prefix(...[]byte)`：按前缀匹配
- `HTTP1Fast()`：用 method 前缀快速判断 HTTP/1（快，但不解析）
- `HTTP1()`：解析 HTTP/1 request（更准）
- `HTTP2()`：匹配 HTTP/2 client preface
- `HTTP2HeaderField(name, value)`：匹配 HTTP/2 HEADERS 中某个字段
- `HTTP2HeaderFieldPrefix(name, valuePrefix)`：匹配 HTTP/2 HEADERS 字段前缀
- `HTTP2HeaderFieldSendSettings(name, value)`：用于 `MatchWithWriters`；sniff 阶段必要时写 SETTINGS，再去匹配 header

## 常见坑 / 设计约束

- **匹配只发生一次**：连接在 Accept 后决定归属，后续不能在同一连接上“切协议”。
- **资源控制**：匹配依赖读取并缓存字节；建议设置 `WithReadTimeout` 和 `WithMaxSniffBytes`。
- **Java gRPC**：部分 Java gRPC 客户端会等服务端 SETTINGS；请用 `MatchWithWriters(HTTP2HeaderFieldSendSettings(...))`。
- **TLS**：如果你在 connmux 之后再做 TLS（或让 `net/http` 依赖对 `net.Conn` 的类型断言识别 TLS），包装连接可能影响某些 TLS 相关识别；如需 `Request.TLS` 等状态，建议先终止 TLS 再做分流。

## 从 soheilhy/cmux 迁移

如果你原来写的是：

```go
grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
```

可以直接改为：

```go
grpcL := m.Match(connmux.HTTP2HeaderField("content-type", "application/grpc"))
```

或者短期内先用兼容包装层：

```go
import "github.com/pubgo/funk/v2/cmux"
```

> 兼容层会把调用转发到 `connmux`。

## 测试

```zsh
go test ./connmux
```
