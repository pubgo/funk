package errorx

import "github.com/rs/zerolog"

import _ "google.golang.org/grpc/codes"
import _ "google.golang.org/grpc/status"

type ss interface {
	Status() int
	ID() string
}

func init() {
	var log = zerolog.Logger{}
	log.Info()
	log.Panic()
	log.Level()
	log.Fatal()
}
