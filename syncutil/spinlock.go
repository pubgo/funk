package syncutil

import (
	"runtime"
	"sync/atomic"
)

type SpinLock struct {
	lock uint32
}

// Lock ...
func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}

// TryLock ...
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

// Unlock ...
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}
