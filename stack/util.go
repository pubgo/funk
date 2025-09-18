package stack

import (
	"runtime/debug"
	
	"github.com/pubgo/funk/monster"
)

var EnablePrintStack = monster.Bool("stack.enable_print_stack", false, "stack enable print stack data")

func PrintStack() {
	if !EnablePrintStack.Get() {
		return
	}

	debug.PrintStack()
}
