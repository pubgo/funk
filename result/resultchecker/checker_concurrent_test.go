package resultchecker_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/result/resultchecker"
)

func TestRegisterErrCheckConcurrent(t *testing.T) {
	t.Cleanup(func() {
		resultchecker.RemoveErrCheck(checkA)
		resultchecker.RemoveErrCheck(checkB)
	})

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				resultchecker.RegisterErrCheck(checkA)
			} else {
				resultchecker.RegisterErrCheck(checkB)
			}
			_ = resultchecker.GetErrChecks()
		}(i)
	}
	wg.Wait()

	checks := resultchecker.GetErrChecks()
	assert.GreaterOrEqual(t, len(checks), 1)
}

func checkA(context.Context, error) error { return nil }
func checkB(context.Context, error) error { return nil }
