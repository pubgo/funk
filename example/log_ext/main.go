package main

import (
	"errors"
	"github.com/pubgo/funk/log"
)

func main() {
	log.New("dd").Info("hello")
	log.New("dd").Error(errors.New("hello"), "test")
}
