package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alecthomas/repr"
	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/pretty"
)

func (e Errors) Append(err error) Errors {
	return append(e, err)
}

func (e Errors) Format(f fmt.State, verb rune) {
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

func (e Errors) String() string {
	var buf = bytes.NewBuffer(nil)

	for i := range e {
		buf.WriteString("====================================================================\n")
		if _err, ok := e[i].(fmt.Stringer); ok {
			buf.WriteString(_err.String())
		} else {
			buf.WriteString(pretty.Sprintln(e[i]))
		}
	}

	return buf.String()
}

func (e Errors) MarshalJSON() ([]byte, error) {
	var errs []interface{}
	for i := range e {
		if _err, ok := e[i].(json.Marshaler); ok {
			errs = append(errs, _err)
		} else if e[i] != nil {
			errs = append(errs, json.RawMessage(repr.String(e[i])))
		}
	}
	return jjson.Marshal(errs)
}

func (e Errors) Error() string {
	if len(e) > 0 {
		return e[0].Error()
	}
	return ""
}

func (e Errors) Unwrap() error {
	if len(e) == 1 {
		return nil
	}

	return e[1:]
}

func (e Errors) As(target interface{}) bool {
	if len(e) > 0 {
		return errors.As(e[0], target)
	}
	return false
}

func (e Errors) Is(target error) bool {
	if len(e) > 0 {
		return errors.Is(e[0], target)
	}
	return false
}

func (e Errors) ErrorOrNil() error {
	if e == nil || len(e) == 0 {
		return nil
	}

	return e
}
