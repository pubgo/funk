package async

import (
	"context"

	"github.com/pubgo/funk/result"
	"go.uber.org/atomic"
)

func NewThread(h func(ctx context.Context) result.Error) *Thread {
	return &Thread{h: h}
}

type Thread struct {
	h      func(ctx context.Context) result.Error
	start  atomic.Bool
	cancel context.CancelFunc
}

func (t *Thread) Start() error {
	if t.start.Load() {
		return nil
	}

	t.cancel = GoCtx(t.h)
	t.start.Store(true)
	return nil
}

func (t *Thread) Stop() error {
	if t.cancel == nil {
		return nil
	}

	t.cancel()
	t.start.Store(false)
	return nil
}
