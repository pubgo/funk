package stack

import (
	"runtime/debug"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/features"
)

// EnablePrintStack enable print stack trace
var EnablePrintStack = features.Bool("stack.enable_print_stack", false, "stack enable print stack trace")

// Print print stack trace
func Print(forces ...bool) {
	if !EnablePrintStack.Value() && !lo.FirstOrEmpty(forces) {
		return
	}

	debug.PrintStack()
}
