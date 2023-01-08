package stack

import (
	"testing"

	"github.com/k0kubun/pp/v3"
)

func TestTrace(t *testing.T) {
	traces := Trace()
	t.Log(pp.Sprint(traces))
}
