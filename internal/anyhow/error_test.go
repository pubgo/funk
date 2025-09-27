package anyhow

import (
	"testing"

	"github.com/pubgo/funk/v2/errors"
)

func TestErrorLog(t *testing.T) {
	ErrOf(errors.New("test")).RecordLog()
}
