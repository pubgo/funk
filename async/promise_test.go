package async

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromise(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		future := Promise(func(resolve func(string), reject func(err error)) {
			time.Sleep(time.Millisecond * 10)
			resolve("ok")
		})

		assert.Equal(t, future.Await().Unwrap(), "ok")
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("err test")
		future := Promise(func(resolve func(string), reject func(err error)) {
			time.Sleep(time.Millisecond * 10)
			reject(err)
		})

		require.ErrorIs(t, future.Await().GetErr(), err)
	})
}

func TestYield(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		iter := Yield(func(yield func(int)) error {
			yield(1)
			yield(2)
			yield(3)
			return nil
		})
		assert.Equal(t, iter.Await().Unwrap(), []int{1, 2, 3})
	})

	err := fmt.Errorf("test error")
	t.Run("err", func(t *testing.T) {
		iter := Yield(func(yield func(int)) error {
			yield(1)
			yield(2)
			yield(3)
			return err
		})
		require.ErrorIs(t, iter.Await().GetErr(), err)
	})
}

func TestGroup(t *testing.T) {
	now := time.Now()
	defer func() {
		t.Log(time.Since(now))
	}()

	rsp := httpGetList().Await()
	assert.NoError(t, rsp.GetErr())
	assert.Equal(t, len(rsp.Unwrap()), 10)

	data := rsp.Unwrap()
	sort.Ints(data)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, data)
}

func TestGroupAwaitOnWorkerError(t *testing.T) {
	err := fmt.Errorf("worker failed")
	iter := Group(func(async func(func() (int, error))) error {
		async(func() (int, error) {
			time.Sleep(20 * time.Millisecond)
			return 1, nil
		})
		async(func() (int, error) {
			time.Sleep(5 * time.Millisecond)
			return 0, err
		})
		return nil
	})

	require.ErrorIs(t, iter.Await().GetErr(), err)
}

func httpGetList() *Iterator[int] {
	return Group(func(async func(func() (int, error))) error {
		for i := 10; i > 0; i-- {
			i := i
			async(func() (int, error) {
				time.Sleep(time.Millisecond * 10)
				return i, nil
			})
		}
		return nil
	})
}
