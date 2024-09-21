package pyroscope

import (
	"net"
	"strings"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/running"
	"github.com/pubgo/lava/core/lifecycle"
	"github.com/samber/lo"
)

func NewClient(c *Config) lifecycle.Handler {
	return func(lc lifecycle.Lifecycle) {
		err := checkPort(c.ServerAddress)
		if err != nil {
			log.Err(err).Msg("pyroscope disabled")
			return
		}

		log.Info().Str("addr", c.ServerAddress).Msg("pyroscope enabled")
		config := pyroscope.Config{
			ApplicationName: running.Project,
			ServerAddress:   c.ServerAddress,
			Logger:          pyroscope.StandardLogger,
			ProfileTypes: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
				pyroscope.ProfileBlockCount,
				pyroscope.ProfileBlockDuration,
			},
			Tags: map[string]string{
				"instance_id": running.InstanceID,
				"version":     running.Version,
				"env":         running.Env,
			},
		}

		pp := assert.Must1(pyroscope.Start(config))
		lc.BeforeStop(func() {
			log.Err(pp.Stop()).Msg("pyroscope stop")
		})
	}
}

func checkPort(address string) error {
	timeout := 3 * time.Second

	conn, err := net.DialTimeout("tcp", lo.LastOrEmpty(strings.Split(address, "//")), timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
