package result

import (
	"fmt"
	"testing"
	"time"
)

func TestYield(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			panic("d")
			return nil
		}).ToResult())

		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			return fmt.Errorf("err test")
		}).ToResult())

		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			return nil
		}).ToResult())
	})
}
