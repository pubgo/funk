package errcheck

import (
	"fmt"
	
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
)

func errMust(err error, args ...any) {
	if generic.IsNil(err) {
		return
	}

	if len(args) > 0 {
		err = errors.Wrap(err, fmt.Sprint(args...))
	}

	err = errors.WrapStack(err)
	errors.Debug(err)
	panic(err)
}
