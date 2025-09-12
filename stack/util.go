package stack

import (
	"runtime/debug"

	"github.com/pubgo/funk/vars"
)

var EnablePrintStack = vars.Bool("stack.enable_print_stack")

func PrintStack() {
	if !EnablePrintStack.Load() {
		return
	}

	debug.PrintStack()
}
