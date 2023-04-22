package errors

import (
	"bytes"
	"fmt"

	jjson "github.com/goccy/go-json"

	"github.com/pubgo/funk/errors/internal"
)

var _ Errors = (*errorsImpl)(nil)

type errorsImpl struct {
	errs []error
}

func (e *errorsImpl) Format(f fmt.State, verb rune) { strFormat(f, verb, e) }
func (e *errorsImpl) Kind() string                  { return "multi" }
func (e *errorsImpl) Errors() []error               { return e.errs }
func (e *errorsImpl) Append(err error) error {
	e.errs = append(e.errs, err)
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
	return jjson.Marshal(errs)
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