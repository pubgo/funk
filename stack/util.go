package stack

import (
	"runtime/debug"

	"github.com/pubgo/funk/v2/features"
)

var EnablePrintStack = features.Bool("stack.enable_print_stack", false, "stack enable print stack data")

func PrintStack() {
	if !EnablePrintStack.GetValue() {
		return
	}

	debug.PrintStack()
}
