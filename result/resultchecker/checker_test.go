package resultchecker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func errCheck1(ctx context.Context, err error) error {
	return err
}

func TestErrCheck(t *testing.T) {
	assert.Equal(t, len(GetErrCheckStacks()), 0)

	assert.Equal(t, RegisterErrCheck(errCheck1), true)
	assert.Equal(t, len(GetErrCheckStacks()), 1)

	assert.Equal(t, RegisterErrCheck(errCheck1), false)
	assert.Equal(t, len(GetErrCheckStacks()), 1)
	assert.Equal(t, GetErrCheckStacks()[0].Short(), "resultchecker/checker_test.go:10 errCheck1")

	RemoveErrCheck(errCheck1)
	assert.Equal(t, len(GetErrCheckStacks()), 0)
}
