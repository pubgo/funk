package strutil

import (
	"fmt"
	"testing"
	"time"
)

func TestFlatten(t *testing.T) {
	t.Log(Flatten(
		"err",
		fmt.Errorf("hello"),
		"hello", "world",
		"hello1", "world world",
		"hello2", 2,
		"hello3", time.Second,
		"hello4", []string{"hello"},
		"hello5", map[string]any{"hello": "hello"},
	))
}
