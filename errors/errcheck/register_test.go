package errcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func errCheck1(ctx context.Context, err error) error {
	return err
}

func TestErrCheck(t *testing.T) {
	assert.Equal(t, len(GetErrCheckFrames()), 0)

	RegisterErrCheck(errCheck1)
	assert.Equal(t, len(GetErrCheckFrames()), 1)

	RegisterErrCheck(errCheck1)
	assert.Equal(t, len(GetErrCheckFrames()), 1)

	RemoveErrCheck(errCheck1)
	assert.Equal(t, len(GetErrCheckFrames()), 0)
}
