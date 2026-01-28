package retry

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/stack"
)

const defaultRetryCount = 3

func MustDo(b Backoff, f func(i int) error) {
	assert.Must(Do(b, f))
}

func Do(b Backoff, f func(i int) error) (err error) {
	wrap := func(i int) (err error) {
		defer recovery.Err(&err)
		return f(i)
	}

	var logger = slog.With("caller", stack.CallerWithFunc(f))

	for i := 0; ; i++ {
		if err = wrap(i); err == nil {
			return nil
		}

		logger.Debug("attempt retry", "count", i, "err", err)
		dur, stop := b.Next()
		if !stop {
			time.Sleep(dur)
			continue
		}

		return fmt.Errorf("retry failed, count %d: %w", i, err)
	}
}

func MustDoVal[T any](b Backoff, f func(i int) (T, error)) T {
	return assert.Must1(DoVal(b, f))
}

func DoVal[T any](b Backoff, f func(i int) (T, error)) (val T, err error) {
	wrap := func(i int) (val T, err error) {
		defer recovery.Err(&err)
		return f(i)
	}

	var logger = slog.With("caller", stack.CallerWithFunc(f))
	for i := 0; ; i++ {
		if val, err = wrap(i); err == nil {
			return val, nil
		}

		logger.Debug("attempt retry", "count", i, "err", err)
		dur, stop := b.Next()
		if !stop {
			time.Sleep(dur)
			continue
		}

		return val, fmt.Errorf("retry failed, count %d: %w", i, err)
	}
}

func Default() Backoff {
	return WithMaxRetries(defaultRetryCount, NewConstant(time.Millisecond*100))
}
