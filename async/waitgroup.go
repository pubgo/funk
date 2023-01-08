package async

import (
	"runtime"
	"sync"
	"sync/atomic"
	_ "unsafe"

	"github.com/pubgo/funk/fastrand"
)

//go:linkname state sync.(*WaitGroup).state
func state(*sync.WaitGroup) (*uint64, *uint32)

var defaultConcurrent = uint32(runtime.NumCPU() * 2)

func NewWaitGroup(maxConcurrent ...uint32) *WaitGroup {
	var c = defaultConcurrent
	if len(maxConcurrent) > 0 {
		c = maxConcurrent[0]
	}
	return &WaitGroup{maxConcurrent: c}
}

type WaitGroup struct {
	wg            sync.WaitGroup
	maxConcurrent uint32
}

func (t *WaitGroup) Count() uint32 {
	count, _ := state(&t.wg)
	return uint32(atomic.LoadUint64(count) >> 32)
}

func (t *WaitGroup) check() {
	// 阻塞, 等待任务处理完毕
	// 采样率(10), 打印log
	if t.Count() >= t.maxConcurrent && fastrand.Sampling(10) {
		logs.Warn().
			Uint32("current", t.Count()).
			Uint32("maximum", t.maxConcurrent).
			Msg("WaitGroup current concurrent number exceeds the maximum concurrent number of the system")
	}
}

func (t *WaitGroup) Inc()          { t.check(); t.wg.Add(1) }
func (t *WaitGroup) Dec()          { t.wg.Done() }
func (t *WaitGroup) Done()         { t.wg.Done() }
func (t *WaitGroup) Wait()         { t.wg.Wait() }
func (t *WaitGroup) Add(delta int) { t.check(); t.wg.Add(delta) }
