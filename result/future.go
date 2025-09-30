package result

import (
	"context"

	"github.com/pubgo/funk/v2/errors"
	"github.com/samber/lo"
)

func AsyncErr(fn func() Error) ErrFuture {
	if fn == nil {
		return ErrFuture{e: errors.WrapCaller(errFnIsNil, 1)}
	}

	var future = newErrFuture()
	go func() { defer future.close(); future.setErr(try(func() error { return fn().getErr() })) }()
	return future
}

func Async[T any](fn func() Result[T]) Future[T] {
	if fn == nil {
		return Future[T]{v: Fail[T](errors.WrapCaller(errFnIsNil, 1))}
	}

	var future = newFuture[T]()
	go func() { defer future.close(); future.setVal(tryResult(fn)) }()
	return future
}

func newFuture[T any]() Future[T] {
	return Future[T]{done: make(chan struct{})}
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

func newErrFuture() ErrFuture {
	return ErrFuture{done: make(chan struct{})}
}

type ErrFuture struct {
	e    error
	done chan struct{}
}

func (f *ErrFuture) close()           { close(f.done) }
func (f *ErrFuture) setErr(err error) { f.e = err }

func (f *ErrFuture) Await(ctxL ...context.Context) Error {
	ctx := lo.FirstOr(ctxL, context.Background())
	select {
	case <-f.done:
		return ErrOf(f.e)
	case <-ctx.Done():
		return ErrOf(ctx.Err())
	}
}
