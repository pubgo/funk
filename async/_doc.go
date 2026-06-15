package async

// Promise/Future, Group/Yield iterators, and Go helpers for async work.
//
// Iterator.Await drains the value channel to completion before returning the
// first error, so callers never leave Group workers blocked on send.
//
// https://github.com/octu0/chanque
