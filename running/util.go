package running

import (
	"fmt"
	rt "runtime"

	semver "github.com/hashicorp/go-version"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo/version"
	"github.com/pubgo/funk/v2/recovery"
)

func GetSysInfo() map[string]string {
	return map[string]string{
		"main_path":     version.MainPath(),
		"grpc_port":     fmt.Sprintf("%v", GrpcPort()),
		"http_post":     fmt.Sprintf("%v", HttpPort()),
		"debug":         fmt.Sprintf("%v", Debug()),
		"cur_dir":       Pwd,
		"local_ip":      LocalIP,
		"namespace":     Namespace,
		"instance_id":   InstanceID,
		"device_id":     DeviceID,
		"project":       Project(),
		"hostname":      Hostname,
		"build_time":    version.BuildTime(),
		"version":       Version(),
		"domain":        Domain(),
		"commit_id":     CommitID(),
		"go_root":       rt.GOROOT(),
		"go_arch":       rt.GOARCH,
		"go_os":         rt.GOOS,
		"go_version":    rt.Version(),
		"num_cpu":       fmt.Sprintf("%v", rt.NumCPU()),
		"num_goroutine": fmt.Sprintf("%v", rt.NumGoroutine()),
	}
}

func CheckVersion() {
	defer recovery.Exit()
	assert.If(version.Project() == "", "project is null")
	assert.If(version.Version() == "", "version is null")
	assert.If(version.CommitID() == "", "commitID is null")
	assert.If(version.BuildTime() == "", "buildTime is null")
	assert.MustFn(func() error {
		_, err := semver.NewVersion(version.Version())
		if err != nil {
			return fmt.Errorf("version(%s) error: %w", version.Version(), err)
		}
		return nil
	})
}
