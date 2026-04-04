package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/internal/errors/errcolorfield"
	"github.com/pubgo/funk/v2/stack"
)

var debugPretty = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(true)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(5)
	return printer
})

func ErrStringify(buf *bytes.Buffer, err error) {
	if err == nil {
		return
	}

	err1, ok := err.(fmt.Stringer)
	if ok {
		if _, ok = err.(*ErrWrap); !ok {
			buf.WriteString("===============================================================\n")
		}
		buf.WriteString(err1.String())
		return
	}

	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorErrMsg, strings.TrimSpace(err.Error()))
	fmt.Fprintf(buf, "%s]: %s\n", errcolorfield.ColorErrDetail, strings.TrimSpace(fmt.Sprintf("%v", err)))
	ErrStringify(buf, Unwrap(err))
}

func ErrJsonify(err error) map[string]any {
	if err == nil {
		return make(map[string]any)
	}

	data := make(map[string]any, 6)
	if _err, ok := err.(json.Marshaler); ok {
		data["cause"] = _err
	} else {
		data["err_msg"] = err.Error()
		data["err_detail"] = fmt.Sprintf("%v", err)
		err = Unwrap(err)
		if err != nil {
			data["cause"] = ErrJsonify(err)
		}
	}
	return data
}

func PrintFormat(f fmt.State, verb rune, err Error) {
	switch verb {
	case 'v':
		data, err := err.MarshalJSON()
		if err != nil {
			_, _ = fmt.Fprintln(f, err.Error())
		} else {
			_, _ = fmt.Fprintln(f, string(data))
		}
	case 's', 'q':
		_, _ = fmt.Fprintln(f, err.String())
	}
}

func getStack() []*stack.Frame {
	return lo.Filter(stack.Trace(), func(item *stack.Frame, index int) bool { return !item.IsRuntime() })
}

func NewErrorId() string { return xid.New().String() }

func getErrorId(err error) string {
	if err == nil {
		return ""
	}

	errId, ok := lo.ErrorsAs[ErrorID](err)
	if ok {
		return errId.ID()
	}

	return NewErrorId()
}
