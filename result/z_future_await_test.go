package result_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgo/funk/v2/result"
)

func TestFutureAwaitPrefersCompletedValueOverCancelledContext(t *testing.T) {
	future := result.Async(func() result.Result[int] {
		return result.OK(42)
	})

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := future.Await(ctx)
	val, ok := r.TryUnwrap()
	require.True(t, ok, "completed work should win over cancelled context, got err=%v", r.Err())
	assert.Equal(t, 42, val)
}

func TestFutureAwaitReturnsTimeoutWhenWorkIsSlow(t *testing.T) {
	future := result.Async(func() result.Result[int] {
		time.Sleep(200 * time.Millisecond)
		return result.OK(42)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	r := future.Await(ctx)
	assert.True(t, r.IsErr())
	assert.ErrorIs(t, r.Err(), context.DeadlineExceeded)
}

func TestErrFutureAwaitReturnsTimeoutWhenWorkIsSlow(t *testing.T) {
	future := result.AsyncErr(func() result.Error {
		time.Sleep(200 * time.Millisecond)
		return result.ErrOf(nil)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	r := future.Await(ctx)
	assert.True(t, r.IsErr())
	assert.ErrorIs(t, r.Err(), context.DeadlineExceeded)
}
