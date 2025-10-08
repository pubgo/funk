package running

import (
	"fmt"
	rt "runtime"

	semver "github.com/hashicorp/go-version"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/buildinfo"
	"github.com/pubgo/funk/v2/recovery"
)

func GetSysInfo() map[string]string {
	return map[string]string{
		"main_path":     buildinfo.MainPath(),
		"grpc_port":     fmt.Sprintf("%v", GrpcPort),
		"http_post":     fmt.Sprintf("%v", HttpPort),
		"debug":         fmt.Sprintf("%v", Debug()),
		"cur_dir":       Pwd,
		"local_ip":      LocalIP,
		"namespace":     Namespace,
		"instance_id":   InstanceID,
		"device_id":     DeviceID,
		"project":       Project,
		"hostname":      Hostname,
		"build_time":    buildinfo.BuildTime(),
		"version":       Version,
		"domain":        Domain,
		"commit_id":     CommitID,
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
	assert.If(buildinfo.Project() == "", "project is null")
	assert.If(buildinfo.Version() == "", "version is null")
	assert.If(buildinfo.CommitID() == "", "commitID is null")
	assert.If(buildinfo.BuildTime() == "", "buildTime is null")
	assert.MustFn(func() error {
		_, err := semver.NewVersion(buildinfo.Version())
		if err != nil {
			return fmt.Errorf("version(%s) error: %w", buildinfo.Version(), err)
		}
		return nil
	})
}
