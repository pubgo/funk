package syncutil

import (
	"runtime"
	"sync"
	"sync/atomic"
	_ "unsafe"

	"github.com/pubgo/funk/async"
	"github.com/pubgo/funk/fastrand"
)

//go:linkname state sync.(*WaitGroup).state
func state(*sync.WaitGroup) (*uint64, *uint32)

var defaultConcurrent = uint32(runtime.NumCPU() * 2)

func NewWaitGroup(maxConcurrent ...uint32) *WaitGroup {
	c := defaultConcurrent
	if len(maxConcurrent) > 0 {
		c = maxConcurrent[0]
	}
	return &WaitGroup{maxConcurrent: c}
}

type WaitGroup struct {
	_             noCopy
	wg            sync.WaitGroup
	maxConcurrent uint32
	err           error
}

func (t *WaitGroup) Count() uint32 {
	count, _ := state(&t.wg)
	return uint32(atomic.LoadUint64(count) >> 32)
}

func (t *WaitGroup) check() {
	// 阻塞, 等待任务处理完毕
	// 采样率(1), 打印log
	for t.Count() >= t.maxConcurrent {
		if fastrand.Sampling(1) {
			logs.Warn().
				Uint32("current", t.Count()).
				Uint32("maximum", t.maxConcurrent).
				Msg("WaitGroup current concurrent number exceeds the maximum concurrent number of the system")
		}

		runtime.Gosched()
	}
}

func (t *WaitGroup) Go(fn func()) {
	t.wg.Add(1)
	t.check()
	async.GoSafe(
		func() error { fn(); return nil },
		func(err error) {
			t.err = err
		},
	)
}

func (t *WaitGroup) Wait() error {
	t.wg.Wait()
	return t.err
}

type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
