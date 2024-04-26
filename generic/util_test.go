package generic

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	data := []int{1, 2, 3, 4}
	t.Log(Map(data, func(i int) string {
		return strconv.Itoa(data[i])
	}))
}
