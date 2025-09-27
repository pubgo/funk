package main

import (
	"os"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/pretty"
)

func main() {
	pretty.Println(os.Environ())
	pretty.Println(env.Map())
}
