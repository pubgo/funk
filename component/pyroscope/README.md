# Pyroscope Component

Grafana [Pyroscope](https://grafana.com/docs/pyroscope/latest/) continuous profiling integration for Go services.

Built on [`github.com/grafana/pyroscope-go`](https://github.com/grafana/pyroscope-go).

## Installation

```bash
go get github.com/pubgo/funk/v2/component/pyroscope
```

## Quick Start

```go
import (
    "github.com/pubgo/funk/v2/component/lifecycle"
    "github.com/pubgo/funk/v2/component/pyroscope"
    "github.com/pubgo/funk/v2/log"
)

func main() {
    lc := /* lifecycle from your app */
    client := pyroscope.New(pyroscope.Param{
        Cfg: &pyroscope.Config{
            Enabled:       true,
            ServerAddress: "http://pyroscope:4040",
        },
        Logger: log.GetLogger("app"),
        Lc:     lc,
    })
    _ = client

    pyroscope.TagWrapper(ctx, pyroscope.Labels("handler", "CreateOrder"), func(ctx context.Context) {
        // profiled code
    })
}
```

## Configuration

| Field | Description |
|-------|-------------|
| `enabled` | Turn profiling on/off |
| `application_name` | Pyroscope app name (default: `{project}/{version}`) |
| `server_address` | Pyroscope server URL |
| `basic_auth_user` / `basic_auth_password` | HTTP basic auth (Grafana Cloud) |
| `tenant_id` | Multi-tenant Pyroscope ID |
| `upload_rate` | Profile upload interval (default: `1m`) |
| `profile_types` | CPU, heap, goroutine, etc. (defaults match pyroscope-go) |
| `tags` | Extra labels merged with hostname/env/version/instance_id |
| `disable_gc_runs` | Pass through to pyroscope-go |
| `disable_log` | Silence pyroscope client logs |

Example YAML:

```yaml
pyroscope:
  enabled: true
  server_address: http://pyroscope:4040
  application_name: my-service
  profile_types:
    - cpu
    - alloc_objects
    - inuse_objects
  tags:
    team: platform
```

## Lifecycle

When `Param.Lc` is set, `BeforeStop` calls `profiler.Stop()` to flush remaining profiles on shutdown.

## Profiling Tags

Use `TagWrapper` or `Labels` (re-exported from pyroscope-go) to attach labels to hot paths:

```go
pyroscope.TagWrapper(ctx, pyroscope.Labels("controller", "slow"), func(ctx context.Context) {
    handleSlowPath(ctx)
})
```

## Pull Mode

For pull-based profiling, enable `net/http/pprof` in your HTTP server; no push client is required. See [Pyroscope Go SDK docs](https://grafana.com/docs/pyroscope/latest/configure-client/language-sdks/go_push/).

## Example

```bash
# requires a running Pyroscope server
PYROSCOPE_SERVER=http://localhost:4040 go run ./component/pyroscope/example
```
