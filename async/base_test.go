package async

import (
	"testing"

	"github.com/pubgo/funk/log"
)

func TestLogger(t *testing.T) {
	logs.Debug().Msg("hello")
	logs.WithName("world").Debug().Msg("hello")
	log.Debug().Msg("hello")
}
