package async

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	ret := Async(func() (*http.Response, error) {
		return http.Get(server.URL)
	}).Await()
	assert.NoError(t, ret.GetErr())
	rsp := ret.Unwrap()
	if b := rsp.Body; b != nil {
		defer func() {
			_ = b.Close()
		}()
	}
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
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
