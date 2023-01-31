package retry

import (
	"time"

	"github.com/pubgo/funk/assert"
)

const DefaultConstant = time.Second

// NewConstant creates a new constant backoff using the value t.
func NewConstant(t time.Duration) Backoff {
	assert.If(t <= 0, "[t] must be greater than 0")

	return BackoffFunc(func() (time.Duration, bool) { return t, false })
}
