package syncx

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pubgo/funk/result"
)

func TestAsync(t *testing.T) {
	t.Log(Async(func() result.Result[*http.Response] {
		return result.New(http.Get("https://www.baidu.com"))
	}).Await())
}

func TestYield(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			panic("d")
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

func httpGetList() result.Chan[*http.Response] {
	return AsyncGroup(func(async func(func() result.Result[*http.Response])) error {
		for i := 2; i > 0; i-- {
			async(func() result.Result[*http.Response] {
				return result.New(http.Get("https://www.baidu.com"))
			})
		}

		return nil
	})
}

func TestGoChan(t *testing.T) {
	var now = time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()

	var val1 = Async(func() result.Result[string] {
		time.Sleep(time.Millisecond)
		fmt.Println("2")
		//return WithErr(errors.New("error"))
		return result.OK("hello")
	})

	var val2 = Async(func() result.Result[string] {
		time.Sleep(time.Millisecond)
		fmt.Println("3")
		//return WithErr(errors.New("error"))
		return result.OK("hello")
	})

	fmt.Println(Wait(val1, val2).ToResult())
}
