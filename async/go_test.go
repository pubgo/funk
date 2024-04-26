package async

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	t.Log(Async(func() (*http.Response, error) {
		return http.Get("https://www.baidu.com")
	}).Await())
}

func TestYield(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			panic("d")
		}).ToList())

		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			return fmt.Errorf("err test")
		}).ToList())

		t.Log(Yield(func(yield func(string)) error {
			yield(time.Now().String())
			yield(time.Now().String())
			yield(time.Now().String())
			return nil
		}).ToList())
	})
}

func httpGetList() *Iterator[*http.Response] {
	return Group(func(async func(func() (*http.Response, error))) error {
		for i := 2; i > 0; i-- {
			async(func() (*http.Response, error) {
				return http.Get("https://www.baidu.com")
			})
		}
		return nil
	})
}

func TestGoChan(t *testing.T) {
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()

	val1 := Async(func() (string, error) {
		time.Sleep(time.Millisecond)
		fmt.Println("2")
		// return WithErr(errors.New("error"))
		return "hello", nil
	})

	val2 := Async(func() (string, error) {
		time.Sleep(time.Millisecond)
		fmt.Println("3")
		// return WithErr(errors.New("error"))
		return "hello", nil
	})

	_ = val1
	_ = val2
}
