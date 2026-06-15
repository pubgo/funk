package resultchecker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/result/resultchecker"
)

func TestCreateCtxAndGetCheckersFromCtx(t *testing.T) {
	check := func(context.Context, error) error { return nil }

	ctx := resultchecker.CreateCtx(nil, []resultchecker.ErrChecker{check})
	got := resultchecker.GetCheckersFromCtx(ctx)
	assert.Len(t, got, 1)

	assert.Nil(t, resultchecker.GetCheckersFromCtx(nil))
	assert.Nil(t, resultchecker.GetCheckersFromCtx(context.Background()))
}
