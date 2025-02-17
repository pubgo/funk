package anyhow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func errCheck1(err error) error {
	return err
}

func TestErrCheck(t *testing.T) {
	assert.Equal(t, len(FindErrCheckList()), 0)

	RegisterErrCheck(errCheck1)

	assert.Equal(t, len(FindErrCheckList()), 1)

	RemoveErrCheck(errCheck1)
	assert.Equal(t, len(FindErrCheckList()), 0)
}
