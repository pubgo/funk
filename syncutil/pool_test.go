package syncutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	now := time.Now()
	defer func() {
		var cost = time.Since(now)
		assert.True(t, cost > time.Millisecond*10*2 && cost < time.Millisecond*10*3)
	}()

	var p = NewPool().WithMaxGoroutines(5)
	assert.Equal(t, p.MaxGoroutines(), 5)
	for i := 0; i < 10; i++ {
		p.Go(func() {
			time.Sleep(time.Millisecond * 10)
		})
	}
	assert.NoError(t, p.Wait())
}
