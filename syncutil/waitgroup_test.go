package syncutil

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	var now = time.Now()
	defer func() {
		t.Log(time.Since(now))
	}()
	var wg = NewWaitGroup(10)
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			time.Sleep(time.Millisecond * 10)
		})
	}
	t.Log(wg.Wait())
}
