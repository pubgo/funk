package cloudevent

// Promise/Future, Group/Yield iterators, and Go helpers for async work.
//
// Protobuf extensions live in github.com/pubgo/funk/v2/proto/cloudeventoption.
// Message types live in github.com/pubgo/funk/v2/proto/cloudevent.
//
// Iterator.Await drains the value channel to completion before returning the
// first error, so callers never leave Group workers blocked on send.
//
// https://github.com/octu0/chanque
