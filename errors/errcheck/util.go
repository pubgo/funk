package errcheck

import (
	"fmt"

	"github.com/pubgo/funk/errors"
)

func errMust(err error, args ...any) {
	if err == nil {
		return
	}

	if len(args) > 0 {
		err = errors.Wrap(err, fmt.Sprint(args...))
	}

	err = errors.WrapStack(err)
	errors.Debug(err)
	panic(err)
}
