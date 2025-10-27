package ctxutil_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/ctxutil"
)

func TestClone(t *testing.T) {
	cc, cancel := context.WithCancel(context.TODO())
	oldCtx := context.WithValue(cc, "hello", "hello")

	newCtx, _ := ctxutil.Clone(oldCtx)
	cancel()

	assert.Equal(t, oldCtx.Err(), context.Canceled)
	assert.Equal(t, newCtx.Value("hello"), "hello")
	assert.Equal(t, newCtx.Err(), nil)
}
