package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {
	var gErr error
	var err = ErrProxyOf(&gErr)
	Errorf("test proxy error").Log().CatchErr(&err)
	assert.NotNil(t, gErr)
	assert.NotNil(t, err.GetErr())
	assert.Equal(t, gErr, err.GetErr())
}
