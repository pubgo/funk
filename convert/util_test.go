package convert

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var vv = Map(map[string]int{"a": 1, "b": 100}, func(s int) string { return fmt.Sprintf("%v", s) })
	assert.Equal(t, vv, map[string]string{"a": "1", "b": "100"})
}

func TestMapL(t *testing.T) {
	var vv = MapL([]int{1, 2, 3, 4}, func(s int) string { return fmt.Sprintf("%v", s*2) })
	assert.Equal(t, vv, []string{"2", "4", "6", "8"})
}
