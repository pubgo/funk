package errinter

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
	"google.golang.org/protobuf/proto"
	
	"github.com/pubgo/funk"
	"github.com/pubgo/funk/proto/errorpb"
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

func GetErrorId(err error) string {
	if err == nil {
		return ""
	}

	for err != nil {
		if v, ok := err.(Error); ok {
			return v.ID()
		}

		err = Unwrap(err)
	}

	return ""
}

func Unwrap(err error) error {
	u, ok := err.(ErrUnwrap)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func ParseErrToPb(err error) proto.Message {
	switch err1 := err.(type) {
	case nil:
		return nil
	case ErrorProto:
		return err1.Proto()
	case GRPCStatus:
		return err1.GRPCStatus().Proto()
	case proto.Message:
		return err1
	default:
		return &errorpb.ErrMsg{Msg: err.Error(), Detail: fmt.Sprintf("%v", err)}
	}
}
