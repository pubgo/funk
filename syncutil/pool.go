package syncutil

import (
	"runtime"
)

// NewPool creates a new Pool.
func NewPool() *Pool {
	return &Pool{limiter: make(limiter, uint32(runtime.NumCPU()*2))}
}

// Pool is a pool of goroutines used to execute tasks concurrently.
//
// Tasks are submitted with Go(). Once all your tasks have been submitted, you
// must call Wait() to clean up any spawned goroutines and propagate any
// panics.
//
// Goroutines are started lazily, so creating a new pool is cheap. There will
// never be more goroutines spawned than there are tasks submitted.
//
// The configuration methods (With*) will panic if they are used after calling
// Go() for the first time.
//
// Pool is efficient, but not zero cost. It should not be used for very short
// tasks. Startup and teardown come with an overhead of around 1µs, and each
// task has an overhead of around 300ns.
type Pool struct {
	handle  WaitGroup
	limiter limiter
}

// Go submits a task to be run in the pool. If all goroutines in the pool
// are busy, a call to Go() will block until the task can be started.
func (p *Pool) Go(f func()) {
	select {
	case p.limiter <- struct{}{}:
		// If we are below our limit, spawn a new worker rather
		// than waiting for one to become available.
		p.handle.Go(func() {
			defer func() { <-p.limiter }()
			f()
		})
	}
}

// Wait cleans up spawned goroutines, propagating any panics that were
// raised by a tasks.
func (p *Pool) Wait() error {
	return p.handle.Wait()
}

// MaxGoroutines returns the maximum size of the pool.
func (p *Pool) MaxGoroutines() int {
	return p.limiter.limit()
}

// WithMaxGoroutines limits the number of goroutines in a pool.
// Defaults to unlimited. Panics if n < 1.
func (p *Pool) WithMaxGoroutines(n int) *Pool {
	if n < 1 {
		panic("max goroutines in a pool must be greater than zero")
	}

	return &Pool{limiter: make(limiter, n)}
}

type limiter chan struct{}

func (l limiter) limit() int {
	return cap(l)
}

func (l limiter) release() {
	if l != nil {
		<-l
	}
}
