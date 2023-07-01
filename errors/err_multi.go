package errors

import (
	"bytes"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/stack"
)

func Append(err error, errs ...error) error {
	if err == nil && len(errs) == 0 {
		return nil
	}

	if len(errs) == 0 {
		return &ErrWrap{
			err:    err,
			caller: stack.Caller(1),
		}
	}

	if err == nil {
		return &ErrWrap{
			err:    &errorsImpl{errs: errs},
			caller: stack.Caller(1),
		}
	}

	var errL []error
	switch err1 := err.(type) {
	case Errors:
		return err1.Append(errs...)
	default:
		errL = make([]error, 0, len(errs)+1)
		errL = append(errL, err1)
		errL = append(errL, errs...)
	}

	return &ErrWrap{
		err:    &errorsImpl{errs: errL},
		caller: stack.Caller(1),
	}
}

var _ Errors = (*errorsImpl)(nil)
var _ fmt.Formatter = (*errorsImpl)(nil)

type errorsImpl struct {
	errs []error
}

func (e *errorsImpl) Format(f fmt.State, verb rune) { strFormat(f, verb, e) }
func (e *errorsImpl) Kind() string                  { return "multi" }
func (e *errorsImpl) Errors() []error               { return e.errs }
func (e *errorsImpl) Append(err ...error) error {
	if len(err) == 0 {
		return e
	}

	errL := make([]error, 0, len(e.errs)+len(err))
	errL = append(errL, e.errs...)
	errL = append(errL, err...)
	e.errs = errL
	return e
}

func (e *errorsImpl) String() string {
	if len(e.errs) == 0 {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, e.Kind()))
	for i := len(e.errs) - 1; i >= 0; i-- {
		errStringify(buf, e.errs[i])
	}
	return buf.String()
}

func (e *errorsImpl) MarshalJSON() ([]byte, error) {
	var errs []interface{}
	for i := range e.errs {
		errs = append(errs, errJsonify(e.errs[i]))
	}
	return json.Marshal(errs)
}

func (e *errorsImpl) Error() string {
	if len(e.errs) > 0 {
		return e.errs[0].Error()
	}
	return ""
}

func (e *errorsImpl) Unwrap() error {
	if len(e.errs) == 1 {
		return e.errs[0]
	}
	return &errorsImpl{errs: e.errs[1:]}
}

func (e *errorsImpl) As(target interface{}) bool {
	if len(e.errs) > 0 {
		return As(e.errs[0], target)
	}
	return false
}

func (e *errorsImpl) Is(target error) bool {
	if len(e.errs) > 0 {
		return Is(e.errs[0], target)
	}
	return false
}
