package result_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/result"
)

func TestAsyncNilFunctionReturnsError(t *testing.T) {
	future := result.Async[int](nil)
	r := future.Await()
	assert.True(t, r.IsErr())
}

func TestAsyncErrNilFunctionReturnsError(t *testing.T) {
	future := result.AsyncErr(nil)
	r := future.Await()
	assert.True(t, r.IsErr())
}

func TestFutureAwaitWithNilContextUsesBackground(t *testing.T) {
	future := result.Async(func() result.Result[int] {
		return result.OK(21)
	})

	r := future.Await(nil)
	assert.Equal(t, 21, r.UnwrapOrEmpty())
}
