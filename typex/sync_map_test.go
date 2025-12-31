package typex

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestName(t *testing.T) {
	var sm SyncMap
	sm.Set("a1", 1)
	sm.Set("a2", 2)
	fmt.Println(sm.Has("a1"))

	lo.Must0(sm.Each(func(key string) {
		fmt.Println(key)
	}))

	lo.Must0(sm.Each(func(key string, val int) {
		fmt.Println(key, val)
	}))

	data := make(map[string]int)
	lo.Must0(sm.MapTo(data))
	fmt.Println(data)

	var data1 map[string]int
	lo.Must0(sm.MapTo(&data1))
	fmt.Println(data1)
}
