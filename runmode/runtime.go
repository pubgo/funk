package runmode

import (
	"os"
	"strings"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/utils"
	"github.com/pubgo/funk/version"
)

// 默认的全局配置
var (
	HttpPort = 8080
	GrpcPort = 50051
	Project  = version.Project()

	IsDebug bool

	// InstanceID service id
	InstanceID = version.InstanceID()

	Version = version.Version()

	CommitID = version.CommitID()

	// Pwd 当前目录
	Pwd = assert.Exit1(os.Getwd())

	// Hostname 主机名
	Hostname = utils.FirstFnNotEmpty(
		func() string { return os.Getenv("HOSTNAME") },
		func() string { return assert.Exit1(os.Hostname()) },
	)

	// Namespace K8s命名空间
	Namespace = utils.FirstFnNotEmpty(
		func() string { return os.Getenv("NAMESPACE") },
		func() string { return os.Getenv("POD_NAMESPACE") },
		func() string {
			var file = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
			if !utils.FileExists(file) {
				return ""
			}

			return strings.TrimSpace(string(assert.Exit1(os.ReadFile(file))))
		},
	)
)

func init() {
	env.GetBoolVal(&IsDebug, "enable_debug", "enable_debug_mode")
}
