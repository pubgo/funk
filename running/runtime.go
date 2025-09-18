package running

import (
	"os"
	"strings"

	"github.com/projectdiscovery/machineid"
	"github.com/rs/xid"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/monster"
	"github.com/pubgo/funk/netutil"
	"github.com/pubgo/funk/pathutil"
	"github.com/pubgo/funk/strutil"
	"github.com/pubgo/funk/version"
)

// default global variables
var (
	HttpPort = 8080
	GrpcPort = 50051
	Project  = version.Project()

	Env = "debug"

	// IsDebug
	// Deprecated: use Debug
	IsDebug = true

	// InstanceID service id
	InstanceID = xid.New().String()

	DeviceID = InstanceID

	Version = version.Version()

	CommitID = version.CommitID()

	Pwd = assert.Exit1(os.Getwd())

	// LocalIP the local IP address of the current service
	LocalIP = netutil.GetLocalIP()

	Hostname = strutil.FirstFnNotEmpty(
		func() string { return os.Getenv("HOSTNAME") },
		func() string { return assert.Exit1(os.Hostname()) },
	)

	// Namespace K8s namespace
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

	Domain string

	enableDebug = monster.Bool("debug", false, "enable debug")
)

func Debug() bool {
	return IsDebug || enableDebug.Get()
}

func init() {
	env.GetBoolVal(&IsDebug, "enable_debug", "debug", "dev_mode")
	env.GetVal(&Env, "env", "run_mode", "run_env")

	id, err := machineid.ID()
	if err == nil {
		DeviceID = id
	}
}
