package async

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsync(t *testing.T) {
	ret := Async(func() (*http.Response, error) { //nolint
		return http.Get("https://httpbin.org")
	}).Await()
	assert.NoError(t, ret.Err())
	rsp := ret.Unwrap()
	if b := rsp.Body; b != nil {
		defer b.Close()
	}
	assert.Equal(t, rsp.StatusCode, 200)
}

func TestGoChan(t *testing.T) {
	now := time.Now()
	defer func() {
		fmt.Println("cost:", time.Since(now))
	}()

	val1 := Async(func() (string, error) {
		time.Sleep(time.Millisecond * 10)
		fmt.Println("1")
		return "hello1", nil
	})

	val2 := Async(func() (string, error) {
		time.Sleep(time.Millisecond * 10)
		fmt.Println("2")
		ret := val1.Await().Unwrap()
		return ret + " hello2", nil
	})

	assert.Equal(t, "hello1 hello2", val2.Await().Unwrap())
}
