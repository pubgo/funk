package stack_test

import (
	"testing"

	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
)

func TestTrace(t *testing.T) {
	traces := stack.Trace()
	t.Log(pretty.Sprint(traces))
}
