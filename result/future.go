package result

import (
	"context"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/errors"
)

func AsyncErr(fn func() Error) *FutureErr {
	if fn == nil {
		return &FutureErr{e: errors.WrapCaller(errFnIsNil, 1)}
	}

	future := newErrFuture()
	go func() {
		defer future.close()
		future.setErr(try(func() error { return fn().getErr() }))
	}()
	return future
}

func Async[T any](fn func() Result[T]) *Future[T] {
	if fn == nil {
		return &Future[T]{v: Fail[T](errors.WrapCaller(errFnIsNil, 1))}
	}

	future := newFuture[T]()
	go func() {
		defer future.close()
		future.setVal(tryResult(fn))
	}()
	return future
}

func newFuture[T any]() *Future[T] {
	return &Future[T]{done: make(chan struct{})}
}

type Future[T any] struct {
	v    Result[T]
	done chan struct{}
}

func (f *Future[T]) close()               { close(f.done) }
func (f *Future[T]) setVal(val Result[T]) { f.v = val }

func (f *Future[T]) Await(ctxL ...context.Context) Result[T] {
	ctx := lo.FirstOr(ctxL, context.Background())
	select {
	case <-f.done:
		return f.v
	case <-ctx.Done():
		return f.v.WithErr(ctx.Err())
	}
}

func newErrFuture() *FutureErr {
	return &FutureErr{done: make(chan struct{})}
}

type FutureErr struct {
	e    error
	done chan struct{}
}

func (f *FutureErr) close()           { close(f.done) }
func (f *FutureErr) setErr(err error) { f.e = err }

func (f *FutureErr) Await(ctxL ...context.Context) Error {
	ctx := lo.FirstOr(ctxL, context.Background())
	select {
	case <-f.done:
		return ErrOf(f.e)
	case <-ctx.Done():
		return ErrOf(ctx.Err())
	}
}
