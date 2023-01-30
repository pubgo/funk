package runmode

import (
	"fmt"
	rt "runtime"

	"github.com/pubgo/funk/version"
)

func GetSysInfo() map[string]string {
	return map[string]string{
		"main_path":     version.MainPath(),
		"grpc_port":     fmt.Sprintf("%v", GrpcPort),
		"http_post":     fmt.Sprintf("%v", HttpPort),
		"debug":         fmt.Sprintf("%v", IsDebug),
		"cur_dir":       Pwd,
		"namespace":     Namespace,
		"instance_id":   InstanceID,
		"project":       Project,
		"hostname":      Hostname,
		"build_time":    version.BuildTime(),
		"version":       Version,
		"commit_id":     CommitID,
		"go_root":       rt.GOROOT(),
		"go_arch":       rt.GOARCH,
		"go_os":         rt.GOOS,
		"go_version":    rt.Version(),
		"num_cpu":       fmt.Sprintf("%v", rt.NumCPU()),
		"num_goroutine": fmt.Sprintf("%v", rt.NumGoroutine()),
	}
}
