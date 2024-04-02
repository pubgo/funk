package typex

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var sm SyncMap
	sm.Set("a1", 1)
	sm.Set("a2", 2)
	fmt.Println(sm.Has("a1"))

	sm.Each(func(key string) {
		fmt.Println(key)
	})

	sm.Each(func(key string, val int) {
		fmt.Println(key, val)
	})

	data := make(map[string]int)
	sm.MapTo(data)
	fmt.Println(data)

	var data1 map[string]int
	sm.MapTo(&data1)
	fmt.Println(data1)
}
