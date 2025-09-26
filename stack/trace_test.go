package stack_test

import (
	"testing"

	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/stack"
	"github.com/samber/lo"
)

func TestTrace(t *testing.T) {
	traces := stack.Trace()
	t.Log(pretty.Sprint(traces))

	traces = lo.Filter(traces, func(item *stack.Frame, index int) bool { return !item.IsRuntime() })
	t.Log(pretty.Sprint(lo.Map(traces, func(item *stack.Frame, index int) string { return item.String() })))
}
