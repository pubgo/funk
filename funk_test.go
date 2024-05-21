package funk

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistinct(t *testing.T) {
	data := Distinct([]int{1, 2, 2, 4, 4, 6, 6, 9, 9, 10}, func(t int) int { return t })
	sort.Ints(data)
	assert.Equal(
		t,
		data,
		[]int{1, 2, 4, 6, 9, 10},
	)
}
