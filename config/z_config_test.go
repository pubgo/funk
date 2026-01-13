package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvMap(t *testing.T) {
	envs := LoadEnvMap("./configs/config.yaml")
	assert.NotNil(t, envs["TEST1"])
	assert.NotNil(t, envs["TEST2"])
}
