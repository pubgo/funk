package stack_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/pubgo/funk/stack"
)

func TestTrace(t *testing.T) {
	traces := stack.Trace()
	t.Log(pp.Sprint(traces))
}
