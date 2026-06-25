package result

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/errors"
)

var (
	_ Checkable = new(Error)
	_ ErrSetter = new(Error)
)

func newError(err error) Error {
	return Error{err: err}
}

type Error struct {
	_ [0]func() // disallow ==

	err error
}

func (e Error) MapErr(fn func(err error) error) Error {
	if e.IsOK() {
		return e
	}

	err := fn(e.getErr())
	err = errors.WrapCaller(err, 1)
	return Error{err: err}
}

func (e Error) LogCtx(ctx context.Context, events ...func(e Event)) Error {
	logErr(ctx, 0, e.err, events...)
	return e
}

func (e Error) Log(events ...func(e Event)) Error {
	logErr(context.Background(), 0, e.err, events...)
	return e
}

// TryUnwrap attempts to unwrap the error, returning it and a boolean indicating if it's an error
// This method is useful when you want to safely extract an error from an Error result
// without triggering a panic. It returns the error (or nil if no error)
// and a boolean indicating whether an error was present.
//
// Example:
//
//	if err, ok := errorResult.TryUnwrap(); ok {
//	    fmt.Printf("Error occurred: %v\n", err)
//	} else {
//	    fmt.Println("No error present")
//	}
func (e Error) TryUnwrap() (error, bool) {
	if e.IsOK() {
		return nil, false
	}
	return e.getErr(), true
}

// Match allows pattern matching on the Error, applying the appropriate function
// This method provides a way to handle both success (no error) and error cases
// in a single operation, similar to pattern matching in functional languages.
//
// Example:
//
//	errorResult.Match(
//	    func() { fmt.Println("No error occurred") },
//	    func(err error) { fmt.Printf("Error: %v\n", err) },
//	)
func (e Error) Match(onOk func(), onErr func(error)) {
	if e.IsOK() {
		onOk()
	} else {
		onErr(e.getErr())
	}
}

// MatchWithError allows pattern matching with an error-returning function
// This method is similar to Match, but both handler functions return an error,
// allowing for chaining operations that may themselves produce errors.
//
// Example:
//
//	err := errorResult.MatchWithError(
//	    func() error { return nil },
//	    func(prevErr error) error { return fmt.Errorf("wrapped: %w", prevErr) },
//	)
func (e Error) MatchWithError(onOk func() error, onErr func(error) error) error {
	if e.IsOK() {
		return onOk()
	}
	return onErr(e.getErr())
}

func (e Error) WithFn(fn func() error) Error {
	if fn == nil {
		return Error{err: errors.WrapCaller(errFnIsNil, 1)}
	}

	err := fn()
	if err == nil {
		return Error{}
	}
	return Error{err: errors.WrapCaller(err, 1)}
}

func (e Error) WithErr(err error, tags ...errors.Tags) Error {
	if err == nil {
		return e
	}

	return Error{err: errors.WrapTagsCaller(err, lo.FirstOrEmpty(tags), 1)}
}

func (e Error) WithErrorf(format string, args ...any) Error {
	return Error{err: errors.WrapCaller(fmt.Errorf(format, args...), 1)}
}

func (e Error) IfErr(fn func(error)) Error {
	if e.IsErr() {
		err := e.getErr()
		fn(err)
	}

	return e
}

func (e Error) InspectErr(fn func(error)) {
	if e.IsErr() {
		fn(e.getErr())
	}
}

func (e Error) Expect(format string, args ...any) {
	if e.IsErr() {
		err := errors.WrapCaller(e.getErr(), 1)
		panicIfError(err, func(e Event) {
			e.Msgf(format, args...)
		})
	}
}

func (e Error) IsErr() bool { return e.getErr() != nil }

func (e Error) IsOK() bool { return e.getErr() == nil }

func (e Error) Err() error { return e.getErr() }

func (e Error) GetErr() error { return e.getErr() }

// Message returns the full human-readable error chain.
func (e Error) Message() string {
	if e.IsOK() {
		return ""
	}
	return errors.FormatChain(e.err)
}

// Tags returns user tags collected from the error chain.
func (e Error) Tags() errors.Tags {
	if e.IsOK() {
		return nil
	}
	return errors.CollectUserTags(e.err)
}

func (e Error) Unwrap() (void Void) {
	if e.IsErr() {
		panicIfError(errors.WrapCaller(e.getErr(), 1))
	}

	return void
}

func (e Error) Must() {
	if e.IsOK() {
		return
	}

	panicIfError(errors.WrapCaller(e.getErr(), 1))
}

func (e Error) ThrowErr(err *error, contexts ...context.Context) bool {
	return catchErr(e, nil, err, contexts...)
}

func (e Error) Throw(setter ErrSetter, contexts ...context.Context) bool {
	return catchErr(e, setter, nil, contexts...)
}

func (e Error) MustWithLog(events ...func(e Event)) {
	if e.IsOK() {
		return
	}

	panicIfError(errors.WrapCaller(e.getErr(), 1), events...)
}

func (e Error) String() string {
	if e.IsOK() {
		return "OK"
	}

	return fmt.Sprintf("Error(%s)", errors.FormatChain(e.err))
}

func (e Error) MarshalJSON() ([]byte, error) {
	if e.IsOK() {
		return []byte("null"), nil
	}

	return errors.JsonPrint(e.err), nil
}

func (e *Error) applyErr(err error) {
	if e == nil || err == nil {
		return
	}
	e.err = err
}

func (e Error) getErr() error { return e.err }

func (e Error) setErrorInner() {
}
