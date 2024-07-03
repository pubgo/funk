package syncutil

import (
	"sync"

	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/stack"
	"github.com/pubgo/funk/try"
	"github.com/rs/zerolog"
)

type WaitGroup struct {
	wg  sync.WaitGroup
	err error
}

func (t *WaitGroup) Go(fn func()) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		err := try.Try(func() error {
			fn()
			return nil
		})
		if err != nil {
			t.err = err
			log.Err(err).Func(func(e *zerolog.Event) {
				e.Str("fn_stack", stack.CallerWithFunc(fn).String())
			}).Msg("recovery func panic")
		}
	}()
}

func (t *WaitGroup) Wait() error {
	t.wg.Wait()
	return t.err
}
