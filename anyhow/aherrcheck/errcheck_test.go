package aherrcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func errCheck1(err error) error {
	return err
}

func TestErrCheck(t *testing.T) {
	assert.Equal(t, len(GetErrCheckFrames()), 0)

	RegisterErrCheck(errCheck1)

	assert.Equal(t, len(GetErrCheckFrames()), 1)

	RemoveErrCheck(errCheck1)
	assert.Equal(t, len(GetErrCheckFrames()), 0)
}
