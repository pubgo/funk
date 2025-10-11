package running

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/projectdiscovery/machineid"
	"github.com/rs/xid"
	"github.com/urfave/cli/v3"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/config"
	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/netutil"
	"github.com/pubgo/funk/v2/pathutil"
	"github.com/pubgo/funk/v2/strutil"
)

// default global variables
var (
	Env      = sync.OnceValue(func() string { return EnvFlag.Get().(string) })
	Debug    = sync.OnceValue(func() bool { return DebugFlag.Get().(bool) })
	HttpPort = sync.OnceValue(func() int { return GrpcPortFlag.Get().(int) })
	GrpcPort = sync.OnceValue(func() int { return HttpPortFlag.Get().(int) })
	Project  = version.Project

	// InstanceID service id
	InstanceID = xid.New().String()

	DeviceID = InstanceID

	Version = version.Version

	CommitID = version.CommitID

	Pwd = assert.Exit1(os.Getwd())

	// LocalIP the local IP address of the current service
	LocalIP = netutil.GetLocalIP()

	Hostname = strutil.FirstFnNotEmpty(
		func() string { return env.Get("HOSTNAME") },
		func() string { return assert.Exit1(os.Hostname()) },
	)

	// Namespace K8s namespace
	Namespace = strutil.FirstFnNotEmpty(
		func() string { return env.Get("NAMESPACE") },
		func() string { return env.Get("POD_NAMESPACE") },
		func() string {
			file := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
			if pathutil.IsNotExist(file) {
				return ""
			}

			return strings.TrimSpace(string(assert.Exit1(os.ReadFile(file))))
		},
	)

	Domain = version.Domain

	DebugFlag = cli.BoolFlag{
		Name:     "debug",
		Usage:    "enable debug mode",
		Value:    false,
		Local:    true,
		Category: "system",
		Sources:  cli.EnvVars(env.Key("enable_debug"), env.Key("debug")),
		Action: func(ctx context.Context, command *cli.Command, b bool) error {
			env.Set("enable_debug", strconv.FormatBool(b))
			env.Set("debug", strconv.FormatBool(b))
			return nil
		},
	}

	EnvFlag = cli.StringFlag{
		Name:     "env",
		Usage:    "running env, dev,test,stage,prod",
		Value:    "dev",
		Local:    true,
		Category: "system",
		Sources:  cli.NewValueSourceChain(cli.EnvVar(env.Key("env")), cli.EnvVar(env.Key("run_env"))),
		Action: func(ctx context.Context, command *cli.Command, s string) error {
			env.Set("env", s)
			env.Set("run_env", s)
			return nil
		},
	}

	GrpcPortFlag = cli.IntFlag{
		Name:     "grpc-port",
		Usage:    "service grpc port",
		Local:    true,
		Value:    50051,
		Category: "system",
		Sources:  cli.EnvVars(env.Key("server_grpc_port")),
		Action: func(ctx context.Context, command *cli.Command, i int) error {
			env.Set("server_grpc_port", strconv.Itoa(i))
			return nil
		},
	}

	HttpPortFlag = cli.IntFlag{
		Name:     "http-port",
		Usage:    "service http port",
		Local:    true,
		Value:    8080,
		Category: "system",
		Sources:  cli.EnvVars(env.Key("server_http_port")),
		Action: func(ctx context.Context, command *cli.Command, i int) error {
			env.Set("server_http_port", strconv.Itoa(i))
			return nil
		},
	}

	ConfFlag = cli.StringFlag{
		Name:     "config",
		Aliases:  []string{"c"},
		Usage:    "config path",
		Value:    config.GetConfigPath(),
		Local:    true,
		Category: "system",
		Sources:  cli.EnvVars(env.Key("config_path")),
		Action: func(ctx context.Context, command *cli.Command, s string) error {
			config.SetConfigPath(s)
			env.Set("config_path", s)
			return nil
		},
	}
)

func init() {
	id, err := machineid.ID()
	if err == nil {
		DeviceID = id
	}
}
