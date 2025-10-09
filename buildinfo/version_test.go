package buildinfo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/module"
)

func TestVersion(t *testing.T) {
	ver := "v0.0.8-0.20251008154318-d8a2f764dac7+dirty"
	assert.True(t, module.IsPseudoVersion(ver))

	a, err := module.PseudoVersionTime(ver)
	assert.NoError(t, err)
	assert.Equal(t, "2025-10-08 15:43:18", a.Format(time.DateTime))

	b, err := module.PseudoVersionRev(ver)
	assert.NoError(t, err)
	assert.Equal(t, "d8a2f764dac7", b)

	c, err := module.PseudoVersionBase(ver)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.7+dirty", c)
}
