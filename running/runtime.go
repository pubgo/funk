package running

import (
	"os"
	"strings"

	"github.com/rs/xid"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/netutil"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/strutil"
	"github.com/pubgo/funk/version"
)

// 默认的全局配置
var (
	HttpPort = 8080
	GrpcPort = 50051
	Project  = version.Project()

	Env     = "debug"
	IsDebug = true

	// InstanceID service id
	InstanceID = xid.New().String()

	Version = version.Version()

	CommitID = version.CommitID()

	// Pwd 当前目录
	Pwd = assert.Exit1(os.Getwd())

	// LocalIP 当前服务的本地IP
	LocalIP = netutil.GetLocalIP()

	// Hostname 主机名
	Hostname = strutil.FirstFnNotEmpty(
		func() string { return os.Getenv("HOSTNAME") },
		func() string { return assert.Exit1(os.Hostname()) },
	)

	// Namespace K8s命名空间
	Namespace = strutil.FirstFnNotEmpty(
		func() string { return os.Getenv("NAMESPACE") },
		func() string { return os.Getenv("POD_NAMESPACE") },
		func() string {
			file := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
			if pathutil.IsNotExist(file) {
				return ""
			}

			return strings.TrimSpace(string(assert.Exit1(os.ReadFile(file))))
		},
	)
)

func init() {
	env.GetBoolVal(&IsDebug, "enable_debug", "debug")
	env.GetWith(&Env, "env", "run_mode")
}
