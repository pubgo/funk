package main

import (
	"fmt"

	"github.com/pubgo/funk/log"
)

var dd = log.GetLogger("dd")

func main() {
	demo(dd)
	demo(dd.WithName("abc"))
}

func demo(base log.Logger) {
	l := base.WithName("MyName").WithName("dd").WithFields(log.Map{"user": "you"})
	l.Info().Fields(map[string]any{"val1": 1, "val2": map[string]int{"k": 1}}).Msg("hello")
	l.Err(nil).Fields(map[string]any{"trouble": true, "reasons": []float64{0.1, 0.11, 3.14}}).Msg("uh oh")
	l.Err(fmt.Errorf("an error occurred")).Int("code", -1).Msg("goodbye")
}
