package generic

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type err1 struct {
}

func (e err1) Error() string {
	return ""
}

func TestMap(t *testing.T) {
	data := []int{1, 2, 3, 4}
	t.Log(Map(data, func(i int) string {
		return strconv.Itoa(data[i])
	}))
}

func TestIsNil(t *testing.T) {
	assert.Equal(t, IsNil(err1{}), false)
	assert.Equal(t, IsNil(struct{}{}), false)
	assert.Equal(t, IsNil(nil), true)
	assert.Equal(t, IsNil(any(nil)), true)
	assert.Equal(t, IsNil((*struct{})(nil)), true)
}
