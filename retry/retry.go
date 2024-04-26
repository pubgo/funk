package retry

import (
	"time"

	"github.com/pubgo/funk/recovery"
)

type Retry func() Backoff

func (d Retry) Do(f func(i int) error) (err error) {
	wrap := func(i int) (err error) {
		defer recovery.Err(&err)
		return f(i)
	}

	b := d()
	for i := 0; ; i++ {
		if err = wrap(i); err == nil {
			return nil
		}

		dur, stop := b.Next()
		if stop {
			return
		}

		time.Sleep(dur)
	}
}

func (d Retry) DoVal(f func(i int) (interface{}, error)) (val interface{}, err error) {
	wrap := func(i int) (val interface{}, err error) {
		defer recovery.Err(&err)
		return f(i)
	}

	b := d()
	for i := 0; ; i++ {
		if val, err = wrap(i); err == nil {
			return val, nil
		}

		dur, stop := b.Next()
		if stop {
			return
		}

		time.Sleep(dur)
	}
}

func New(bs ...Backoff) Retry {
	b := WithMaxRetries(3, NewConstant(DefaultConstant))
	if len(bs) > 0 {
		b = bs[0]
	}

	return func() Backoff { return b }
}

func Default() Retry {
	b := WithMaxRetries(3, NewConstant(time.Millisecond*10))
	return func() Backoff { return b }
}
