package config

import (
	"github.com/pubgo/funk/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComponentName(t *testing.T) {
	is := assert.New(t)
	is.Equal(getComponentName(nil), defaultComponentKey)
	is.Equal(getComponentName(map[string]interface{}{componentConfigKey: nil}), defaultComponentKey)
	is.Equal(getComponentName(map[string]interface{}{componentConfigKey: "hello"}), "hello")
}

func TestName(t *testing.T) {
	errors.RegStackFilter()
}
