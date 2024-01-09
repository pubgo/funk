package utils

import (
	"io"

	"github.com/pubgo/funk/log"
)

func SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Warn().Err(err).Msg("Close operation failed")
	}
}
