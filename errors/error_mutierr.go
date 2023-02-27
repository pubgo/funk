package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/errors/internal"
	"github.com/pubgo/funk/pretty"
)

var _ Errors = (*errorsImpl)(nil)

type errorsImpl struct {
	errs []error
}

func (e *errorsImpl) Kind() string {
	return "multi"
}

func (e *errorsImpl) Errors() []error {
	return e.errs
}

func (e *errorsImpl) Append(err error) error {
	e.errs = append(e.errs, err)
	return e
}

func (e *errorsImpl) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		var data, err = e.MarshalJSON()
		if err != nil {
			fmt.Fprintln(f, err.Error())
		} else {
			fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		fmt.Fprintln(f, e.String())
	}
}

func (e *errorsImpl) String() string {
	if len(e.errs) == 0 {
		return ""
	}

	var buf = bytes.NewBuffer(nil)
	for i := range e.errs {
		buf.WriteString("====================================================================\n")
		if _err, ok := e.errs[i].(fmt.Stringer); ok {
			buf.WriteString(_err.String())
		} else {
			buf.WriteString(fmt.Sprintf("%s]: %q\n", internal.ColorKind, e.Kind()))
			buf.WriteString(fmt.Sprintf("%s]: %s\n", internal.ColorErrDetail, pretty.Sprint(e.errs[i])))
		}
	}

	return buf.String()
}

func (e *errorsImpl) MarshalJSON() ([]byte, error) {
	var errs []interface{}
	for i := range e.errs {
		if _err, ok := e.errs[i].(json.Marshaler); ok {
			errs = append(errs, _err)
		} else if e.errs[i] != nil {
			errs = append(errs, json.RawMessage(repr.String(e.errs[i])))
		}
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
		return nil
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
