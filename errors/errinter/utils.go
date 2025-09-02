package errinter

import (
	"errors"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
	"github.com/pubgo/funk"
)

func ParseError(val interface{}) error {
	if funk.IsNil(val) {
		return nil
	}

	switch v := val.(type) {
	case nil:
		return nil
	case error:
		return v
	case string:
		return errors.New(v)
	case []byte:
		return errors.New(string(v))
	default:
		return errors.New(SimplePrint(v))
	}
}

var Simple = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(3)
	return printer
})

func SimplePrint(v interface{}) string {
	return strings.ReplaceAll(Simple().Sprint(v), "\n", "")
}
