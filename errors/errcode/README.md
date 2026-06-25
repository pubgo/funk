# errcode

`errcode` adds typed business error codes with protobuf and gRPC compatibility on top of the core `errors` package.

## Quick Start

```go
import (
    "github.com/pubgo/funk/v2/errors/errcode"
    "github.com/pubgo/funk/v2/proto/errorpb"
)

const notFound = "demo.user.not_found"

func init() {
    errcode.MustRegisterErrCode(&errorpb.ErrCode{
        Code:       404,
        Message:    "user not found",
        Name:       notFound,
        StatusCode: errorpb.Code_NotFound,
    })
}

func loadUser() error {
    code, ok := errcode.LookupErrCode(notFound)
    if !ok {
        return errcode.NewCodeErr(&errorpb.ErrCode{Name: "demo.internal", Message: "missing code"})
    }
    return errcode.NewCodeErr(code)
}
```

## Registry API

| Function | Description |
|----------|-------------|
| `RegisterErrCode(code)` | Register and return an error on duplicate names |
| `MustRegisterErrCode(code)` | Register and panic on failure |
| `LookupErrCode(name)` | Clone a registered code by name |
| `GetErrCodes()` | Return all registered codes |
| `RegisterErrCodes(code)` | Deprecated alias of `MustRegisterErrCode` |

Generated protobuf plugins typically call `RegisterErrCodes` during package init.

## Error Construction

| Function | Description |
|----------|-------------|
| `NewCodeErr(code, details...)` | Create a typed code error |
| `NewCodeErrWithMsg(code, msg, details...)` | Override the message |
| `WrapCode(err, code)` | Attach a code to an existing error |
| `ParseError(err)` | Extract `*errorpb.ErrCode` from any error |

## Protocol Mapping

| Function | Description |
|----------|-------------|
| `ConvertErr2Status(code)` | Build gRPC status from protobuf code |
| `Err2GrpcCode(err)` | Map standard errors to gRPC codes |
| `GrpcCodeToHTTP(code)` | Map gRPC codes to HTTP status |
| `Http2GrpcCode(status)` | Map HTTP status to gRPC codes |

## See Also

- [errors README](../README.md)
- [errorpb definitions](../../proto/errorpb/)
