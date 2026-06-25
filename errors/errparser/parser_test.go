package errparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/errors/errparser"
)

type customValue struct {
	msg string
}

func TestParse(t *testing.T) {
	assert.Nil(t, errparser.Parse(nil))

	orig := errors.New("wrapped")
	assert.Equal(t, orig, errparser.Parse(orig))

	strErr := errparser.Parse("string error")
	assert.Equal(t, "string error", strErr.Error())
	_, ok := errors.AsA[*errors.ErrWrap](strErr)
	assert.True(t, ok)

	bytesErr := errparser.Parse([]byte("bytes error"))
	assert.Equal(t, "bytes error", bytesErr.Error())

	otherErr := errparser.Parse(42)
	assert.Equal(t, "42", otherErr.Error())
}

func TestRegisterParser(t *testing.T) {
	custom := func(val any) (error, bool) {
		if v, ok := val.(customValue); ok {
			return errors.New("parsed " + v.msg), true
		}
		return nil, false
	}

	assert.True(t, errparser.RegisterParser(custom))
	assert.False(t, errparser.RegisterParser(custom))
	assert.False(t, errparser.RegisterParser(nil))

	parsed := errparser.Parse(customValue{msg: "custom"})
	assert.Equal(t, "parsed custom", parsed.Error())
}
