// https://github.com/chebyrash/promise

package typex

type Result[T any] struct {
	err error
	val T
}

func (v Result[T]) IsErr() bool { return v.err != nil }
func (v Result[T]) Get() T      { return v.val }
func (v Result[T]) Err() error  { return v.err }
func (v Result[T]) WithVal(val T) Result[T] {
	v.val = val
	return v
}

func (v Result[T]) WithErr(err error) Result[T] {
	v.err = err
	return v
}

func OK[T any](val T) Result[T] {
	return Result[T]{val: val}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func Wrap[T any](val T, err error) Result[T] {
	return Result[T]{val: val, err: err}
}
