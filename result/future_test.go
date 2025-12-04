package result_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/result"
)

func TestFuture_ConcurrentAwait(t *testing.T) {
	future := result.Async(func() result.Result[int] {
		time.Sleep(10 * time.Millisecond)
		return result.OK(42)
	})

	var wg sync.WaitGroup
	results := make([]result.Result[int], 10)

	// Start multiple goroutines waiting for the same future
	for i := 0; i < 10; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			results[idx] = future.Await()
		}()
	}

	wg.Wait()

	// All results should be the same
	for _, r := range results {
		assert.True(t, r.IsOK())
		assert.Equal(t, 42, r.UnwrapOrEmpty())
	}
}

func TestErrFuture_ConcurrentAwait(t *testing.T) {
	future := result.AsyncErr(func() (r result.Error) {
		time.Sleep(10 * time.Millisecond)
		return r.WithErr(nil)
	})

	var wg sync.WaitGroup
	results := make([]result.Error, 10)

	// Start multiple goroutines waiting for the same future
	for i := 0; i < 10; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			results[idx] = future.Await()
		}()
	}

	wg.Wait()

	// All results should be the same
	for _, r := range results {
		assert.True(t, r.IsOK())
	}
}
