package syncutil

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
	"github.com/rs/zerolog"
)

var defaultConcurrent = uint32(runtime.NumCPU() * 2)

func NewWaitGroup(maxConcurrent ...uint32) *WaitGroup {
	c := defaultConcurrent
	if len(maxConcurrent) > 0 {
		c = maxConcurrent[0]
	}
	return &WaitGroup{maxConcurrent: c}
}

type WaitGroup struct {
	wg            sync.WaitGroup
	count         atomic.Uint32
	maxConcurrent uint32
}

func (t *WaitGroup) Count() uint32 { return t.count.Load() }

func (t *WaitGroup) checkAndWait() {
	for {
		count := t.Count()
		if count <= t.maxConcurrent {
			break
		}

		if count < defaultConcurrent {
			runtime.Gosched()
		} else {
			time.Sleep(time.Microsecond * 10)
		}
	}
}

func (t *WaitGroup) Go(fn func()) {
	t.wg.Add(1)
	t.checkAndWait()
	go func() {
		defer t.wg.Done()
		err := try.Try(func() error {
			fn()
			return nil
		})
		if err != nil {
			log.Err(err).
				Func(func(e *zerolog.Event) {
					e.Str("fn_stack", stack.CallerWithFunc(fn).String())
				}).Msg("recovery func panic")
		}
	}()
}

func (t *WaitGroup) Wait() {
	t.wg.Wait()
}
