package stack_test

import (
	"testing"

	"github.com/samber/lo"

	"github.com/pubgo/funk/v2/pretty"
	"github.com/pubgo/funk/v2/stack"
)

func TestTrace(t *testing.T) {
	traces := stack.Trace()
	t.Log(pretty.Sprint(traces))

	traces = lo.Filter(traces, func(item *stack.Frame, _ int) bool { return !item.IsRuntime() })
	t.Log(pretty.Sprint(lo.Map(traces, func(item *stack.Frame, _ int) string { return item.String() })))
}
