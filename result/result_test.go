package result_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
	"github.com/pubgo/funk/v2/recovery"
	"github.com/pubgo/funk/v2/result"
	"github.com/pubgo/funk/v2/result/resultchecker"
)

func TestMust(t *testing.T) {
	defer recovery.Testing(t)
	assert.Panics(t, func() {
		result.Must(fmt.Errorf("test must"))
	})
}

type hello struct {
	Name string `json:"name"`
}

func TestName(t *testing.T) {
	defer recovery.DebugPrint()
	ok := &hello{Name: "abc"}
	okBytes := result.Wrap(json.Marshal(&ok))
	data := string(okBytes.Expect("failed to encode json data"))
	t.Log(data)
	if data != `{"name":"abc"}` {
		t.Log(data)
		t.Fatal("not match")
	}
}

func TestResultDo(t *testing.T) {
	ok := result.OK(&hello{Name: "abc"})
	ok.IfOK(func(v *hello) {
		assert.Equal(t, "abc", v.Name)
	})
	ok.IfErr(func(err error) {
		t.Log(err)
	})
}

func TestErrOf(t *testing.T) {
	resultchecker.RegisterErrCheck(log.RecordErr())

	fn1().IfErr(func(err error) {
		errors.DebugPrint(err)
	})
}

func fn1() (r result.Result[string]) {
	if fn3().Throw(&r) {
		return r
	}

	// Use UnwrapErr instead of Must
	val, err := fn2().UnwrapErr()
	if err != nil {
		r = result.Fail[string](err)
		return r
	}

	return result.OK(val)
}

func fn2() (r result.Result[string]) {
	fn3().
		IfErr(func(err error) {
			log.Err(err).Msg("test error")
		}).
		Throw(&r)
	if r.IsErr() {
		return r
	}

	return result.OK("ok")
}

func fn3() result.Error {
	return result.ErrOf(fmt.Errorf("error test, this is error")).
		IfErr(func(err error) {
			log.Err(err).Msg("ddd")
		}).
		Log()
}

func TestMoreReasonableErrorHandling(t *testing.T) {
	// Test our new more reasonable error handling approach
	// This demonstrates how to handle errors without panicking

	// Create a successful result
	okResult := result.OK(42)

	// Convert to another type without panicking
	strResult := result.MapTo(okResult, func(i int) string {
		return fmt.Sprintf("number: %d", i)
	})

	// Apply a function that might fail
	finalResult := result.FlatMapTo(strResult, func(s string) result.Result[int] {
		if len(s) > 10 {
			return result.Fail[int](fmt.Errorf("string too long"))
		}
		return result.OK(len(s))
	})

	// Check the result without panicking
	if finalResult.IsErr() {
		t.Logf("Got expected error: %v", finalResult.GetErr())
	} else {
		value := finalResult.UnwrapOr(-1)
		t.Logf("Got value: %d", value)
	}

	// Test with an error result
	errResult := result.Fail[string](fmt.Errorf("initial error"))

	// Chain operations on an error result
	chainedResult := result.FlatMapTo(errResult, func(s string) result.Result[int] {
		return result.OK(100)
	})

	if chainedResult.IsErr() {
		t.Logf("Chained result also has error: %v", chainedResult.GetErr())
	}
}

func TestNewFunctionality(t *testing.T) {
	// Test Try function
	result1 := result.Try(func() int {
		return 42
	})
	if val, ok := result1.TryUnwrap(); ok {
		t.Logf("Try success: %d", val)
	}

	// Test Try with panic
	result2 := result.Try(func() int {
		panic("something went wrong")
	})
	if result2.IsErr() {
		t.Logf("Try caught panic: %v", result2.GetErr())
	}

	// Test Partition
	results := []result.Result[int]{
		result.OK(1),
		result.Fail[int](fmt.Errorf("error 1")),
		result.OK(2),
		result.Fail[int](fmt.Errorf("error 2")),
		result.OK(3),
	}

	values, errors := result.Partition(results)
	t.Logf("Partitioned values: %v, errors: %v", values, errors)

	// Test Collect
	goodResults := []result.Result[int]{
		result.OK(1),
		result.OK(2),
		result.OK(3),
	}

	collected := result.Collect(goodResults)
	if vals, ok := collected.TryUnwrap(); ok {
		t.Logf("Collected values: %v", vals)
	}

	// Test Match
	result.OK(42).Match(
		func(val int) { t.Logf("Matched OK value: %d", val) },
		func(err error) { t.Logf("Matched error: %v", err) },
	)

	failResult := result.Fail[string](fmt.Errorf("test error"))
	failResult.Match(
		func(val string) { t.Logf("Matched OK value: %s", val) },
		func(err error) { t.Logf("Matched error: %v", err) },
	)
}
