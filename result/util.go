package result

import (
	"sync"

	"github.com/pubgo/funk/xtry"
)

func Yield[T any](do func(yield func(T)) error) Chan[T] {
	var dd = make(chan Result[T])
	go func() {
		defer close(dd)
		err := xtry.Try(func() error { return do(func(t T) { dd <- OK(t) }) })
		if err != nil {
			dd <- Err[T](err)
		}
	}()
	return dd
}

func Async[T any](do func(yield func(func() T)) error) Chan[T] {
	var dd = make(chan Result[T])
	go func() {
		var wg sync.WaitGroup
		var handler = func() error {
			return do(func(f func() T) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := xtry.Try(func() error { dd <- OK(f()); return nil }); err != nil {
						dd <- Err[T](err)
					}
				}()
			})
		}

		if err := xtry.Try(handler); err != nil {
			dd <- Err[T](err)
		}
		wg.Wait()
		close(dd)
	}()
	return dd
}
