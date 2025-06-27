package errors

import (
	"github.com/pubgo/funk/stack"
)

var stackFilters []func(frame *stack.Frame) bool

// RegStackPkgFilter filter fn , pkg
func RegStackPkgFilter(filters ...func(frame *stack.Frame) bool) {
	if len(filters) == 0 {
		return
	}

	stackFilters = append(stackFilters, filters...)
}

func filterStack(frame *stack.Frame) bool {
	for _, filter := range stackFilters {
		if filter(frame) {
			return true
		}
	}
	return false
}
