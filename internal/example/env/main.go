package main

import (
	"os"

	"github.com/pubgo/funk/env"
	"github.com/pubgo/funk/pretty"
)

func main() {
	pretty.Println(os.Environ())
	pretty.Println(env.Map())
}
