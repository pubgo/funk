package stack

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/kr/pretty"
)

func TestTrace(t *testing.T) {
	traces := Trace()
	pretty.Log(traces)
	pp.Println(traces)
}
