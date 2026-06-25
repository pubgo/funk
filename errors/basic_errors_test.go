package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type customError struct{ msg string }

func (e *customError) Error() string {
	return e.msg
}

func TestNew(t *testing.T) {
	err := New("test error")
	assert.NotNil(t, err)
	assert.Equal(t, "test error", err.Error())
}

func TestErrorf(t *testing.T) {
	err := Errorf("test error %d", 42)
	assert.NotNil(t, err)
	assert.Equal(t, "test error 42", err.Error())
}

func TestWrap(t *testing.T) {
	origErr := New("original error")
	wrappedErr := Wrap(origErr, "wrapped message")
	assert.NotNil(t, wrappedErr)
	assert.Equal(t, "original error", wrappedErr.Error()) // 底层错误信息不变

	// 测试包装nil错误
	nilErr := Wrap(nil, "should be nil")
	assert.Nil(t, nilErr)

	// 测试多层包装
	wrappedAgain := Wrap(wrappedErr, "wrapped again")
	assert.NotNil(t, wrappedAgain)
	assert.Equal(t, "original error", wrappedAgain.Error())
}

func TestWrapf(t *testing.T) {
	origErr := New("original error")
	wrappedErr := Wrapf(origErr, "wrapped %s", "message")
	assert.NotNil(t, wrappedErr)
	assert.Equal(t, "original error", wrappedErr.Error())

	// 测试包装nil错误
	nilErr := Wrapf(nil, "should be %s", "nil")
	assert.Nil(t, nilErr)

	// 测试格式化参数
	formattedErr := Wrapf(origErr, "error code: %d", 404)
	assert.NotNil(t, formattedErr)
	assert.Equal(t, "original error", formattedErr.Error())
}

func TestWrapTags(t *testing.T) {
	origErr := New("original error")
	tags := Tags{"key1": "value1", "key2": 42}
	wrappedErr := WrapTags(origErr, tags)
	assert.NotNil(t, wrappedErr)
	assert.Equal(t, "original error", wrappedErr.Error())

	// 测试包装nil错误
	nilErr := WrapTags(nil, tags)
	assert.Nil(t, nilErr)

	// 测试空标签
	emptyTagsErr := WrapTags(origErr, Tags{})
	assert.NotNil(t, emptyTagsErr)
	assert.Equal(t, "original error", emptyTagsErr.Error())
}

func TestWrapStack(t *testing.T) {
	origErr := New("original error")
	stackErr := WrapStack(origErr)
	assert.NotNil(t, stackErr)
	assert.Equal(t, "original error", stackErr.Error())

	// 测试包装nil错误
	nilErr := WrapStack(nil)
	assert.Nil(t, nilErr)

	// 测试包含堆栈信息
	// 注意：这里我们只验证不panic，具体的堆栈内容难以预测
	assert.NotPanics(t, func() {
		_ = stackErr.Error()
	})
}

func TestWrapCaller(t *testing.T) {
	origErr := New("original error")
	callerErr := WrapCaller(origErr)
	assert.NotNil(t, callerErr)
	assert.Equal(t, "original error", callerErr.Error())

	// 测试包装nil错误
	nilErr := WrapCaller(nil)
	assert.Nil(t, nilErr)

	// 测试带跳过帧数的调用
	callerErrWithSkip := WrapCaller(origErr, 2)
	assert.NotNil(t, callerErrWithSkip)
	assert.Equal(t, "original error", callerErrWithSkip.Error())
}

func TestIs(t *testing.T) {
	err1 := New("error1")
	err2 := New("error2")
	wrapped := Wrap(err1, "wrapped")

	assert.True(t, Is(wrapped, err1))
	assert.False(t, Is(wrapped, err2))

	// 测试nil情况
	assert.False(t, Is(nil, err1))
	assert.False(t, Is(err1, nil))
	assert.True(t, Is(nil, nil))
}

func TestAs(t *testing.T) {
	err := &customError{msg: "custom error"}
	wrapped := Wrap(err, "wrapped")

	var target *customError
	ok := As(wrapped, &target)
	assert.True(t, ok)
	assert.Equal(t, "custom error", target.msg)

	// 测试不匹配的情况
	var otherTarget *error
	ok = As(wrapped, &otherTarget)
	assert.False(t, ok)

	// 测试nil情况
	ok = As(nil, &target)
	assert.False(t, ok)
}

func TestAsA(t *testing.T) {
	err := &customError{msg: "custom error"}
	wrapped := Wrap(err, "wrapped")

	target, ok := AsA[*customError](wrapped)
	assert.True(t, ok)
	assert.Equal(t, "custom error", (*target).msg)

	// 测试不匹配的情况 - 使用一个不可能匹配的类型
	_, ok = AsA[*bytes.Buffer](wrapped)
	assert.False(t, ok)

	// 测试nil情况
	_, ok = AsA[*customError](nil)
	assert.False(t, ok)
}

func TestJoin(t *testing.T) {
	err1 := New("error1")
	err2 := New("error2")

	joined := Join(err1, err2)
	assert.NotNil(t, joined)

	// 验证组合错误包含所有子错误
	assert.True(t, Is(joined, err1))
	assert.True(t, Is(joined, err2))

	// 测试单个错误
	singleJoined := Join(err1)
	// 当只有一个非nil错误时，Join的行为可能因标准库实现而异，我们只验证不为nil
	assert.NotNil(t, singleJoined)

	// 测试nil错误
	nilJoined := Join(nil, err1, nil, err2)
	assert.NotNil(t, nilJoined)
	assert.True(t, Is(nilJoined, err1))
	assert.True(t, Is(nilJoined, err2))

	// 测试全部nil错误
	allNilJoined := Join(nil, nil, nil)
	assert.Nil(t, allNilJoined)

	// 测试错误文本
	assert.Contains(t, joined.Error(), "error1")
	assert.Contains(t, joined.Error(), "error2")
}

func TestJsonPrint(t *testing.T) {
	err := New("json test error")

	jsonData := JsonPrint(err)
	assert.NotNil(t, jsonData)

	// 验证是有效的JSON
	var result map[string]any
	err2 := json.Unmarshal(jsonData, &result)
	assert.NoError(t, err2)
	// 注意：JsonPrint可能会返回不同的格式，所以我们检查是否包含错误消息
	assert.Contains(t, string(jsonData), "json test error")
}

func TestGetErrorId(t *testing.T) {
	err := New("test error")
	id := GetErrorId(err)
	assert.NotEmpty(t, id)

	// 测试nil错误
	nilId := GetErrorId(nil)
	assert.Empty(t, nilId)
}

func TestTags_ToMapString(t *testing.T) {
	tags := Tags{
		"string_key": "string_value",
		"int_key":    42,
		"bool_key":   true,
	}

	mapString := tags.ToMapString()
	assert.Equal(t, "string_value", mapString["string_key"])
	assert.Equal(t, "42", mapString["int_key"])
	assert.Equal(t, "true", mapString["bool_key"])
}

func TestIfErr(t *testing.T) {
	// 测试无错误的情况
	result := IfErr(nil, func(err error) error {
		return Wrap(err, "should not be called")
	})
	assert.Nil(t, result)

	// 测试有错误的情况
	origErr := New("original error")
	result = IfErr(origErr, func(err error) error {
		return Wrap(err, "wrapped")
	})
	assert.NotNil(t, result)
	assert.True(t, Is(result, origErr))
}

func TestUnwrap(t *testing.T) {
	// 测试可解包的错误
	origErr := New("original error")
	wrapped := Wrap(origErr, "wrapped")
	unwrapped := Unwrap(wrapped)
	assert.Equal(t, origErr, unwrapped)

	// 测试不可解包的错误
	stdErr := fmt.Errorf("standard error")
	unwrappedStd := Unwrap(stdErr)
	assert.Nil(t, unwrappedStd)

	// 测试nil错误
	unwrappedNil := Unwrap(nil)
	assert.Nil(t, unwrappedNil)

	// 测试多层包装的解包
	innerErr := New("inner error")
	middleWrapped := Wrap(innerErr, "middle")
	outerWrapped := Wrap(middleWrapped, "outer")
	unwrappedOnce := Unwrap(outerWrapped)
	assert.Equal(t, middleWrapped, unwrappedOnce)
	unwrappedTwice := Unwrap(unwrappedOnce)
	assert.Equal(t, innerErr, unwrappedTwice)
}

func TestErrStringify(t *testing.T) {
	// 测试基本错误的字符串化
	err := New("test error")
	buf := &bytes.Buffer{}
	ErrStringify(buf, err)
	result := buf.String()
	assert.Contains(t, result, "test error")

	// 测试nil错误
	buf.Reset()
	ErrStringify(buf, nil)
	assert.Empty(t, buf.String())
}

func TestDebugPrint(t *testing.T) {
	// 测试错误的调试打印
	err := New("debug test error")

	// 确保不会panic
	assert.NotPanics(t, func() {
		DebugPrint(err)
	})

	// 测试nil错误
	assert.NotPanics(t, func() {
		DebugPrint(nil)
	})
}

func TestWrapFn(t *testing.T) {
	origErr := New("original error")
	wrappedErr := WrapFn(origErr, func() Tags {
		return Tags{"dynamic": "value"}
	})
	assert.NotNil(t, wrappedErr)
	assert.Equal(t, "original error", wrappedErr.Error())

	// 测试包装nil错误
	nilErr := WrapFn(nil, func() Tags {
		return Tags{"dynamic": "value"}
	})
	assert.Nil(t, nilErr)
}

func TestWrapKV(t *testing.T) {
	origErr := New("original error")
	wrappedErr := WrapKV(origErr, "key", "value")
	assert.NotNil(t, wrappedErr)
	assert.Equal(t, "original error", wrappedErr.Error())

	// 测试包装nil错误
	nilErr := WrapKV(nil, "key", "value")
	assert.Nil(t, nilErr)
}

func TestGetTags(t *testing.T) {
	err := New("payment failed", Tags{"transaction_id": "txn_123"})
	tags := GetTags(err)
	assert.Equal(t, "txn_123", tags["transaction_id"])

	wrapped := WrapTags(err, Tags{"layer": "service"})
	tags = GetTags(wrapped)
	assert.Equal(t, "txn_123", tags["transaction_id"])

	collected := CollectTags(wrapped)
	assert.Equal(t, "txn_123", collected["transaction_id"])
	assert.Equal(t, "service", collected["layer"])
}

func TestCollectTagsNil(t *testing.T) {
	assert.Nil(t, GetTags(nil))
	assert.Nil(t, CollectTags(nil))
}

func TestTagsCloneOnWrap(t *testing.T) {
	tags := Tags{"key": "value"}
	err := WrapTags(New("original"), tags)
	tags["key"] = "mutated"

	wrapped, ok := AsA[*ErrWrap](err)
	assert.True(t, ok)
	assert.Equal(t, "value", (*wrapped).Tags["key"])
}

func TestUnwrapStdlib(t *testing.T) {
	inner := New("inner")
	wrapped := fmt.Errorf("outer: %w", inner)
	assert.Equal(t, inner, Unwrap(wrapped))
}

func TestMarshalError(t *testing.T) {
	data, err := MarshalError(New("marshal test"))
	assert.NoError(t, err)
	assert.Contains(t, string(data), "marshal test")

	data, err = MarshalError(nil)
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestRootCauseAndWalk(t *testing.T) {
	root := &customError{msg: "root"}
	outer := Wrap(Wrap(root, "middle"), "outer")
	assert.Equal(t, root, RootCause(outer))

	var layers int
	Walk(outer, func(err error) bool {
		layers++
		return true
	})
	assert.GreaterOrEqual(t, layers, 3)
}

func TestWrapStackCapturesStacks(t *testing.T) {
	stackErr := WrapStack(New("stack test"))
	wrapped, ok := AsA[*ErrWrap](stackErr)
	assert.True(t, ok)
	assert.NotEmpty(t, (*wrapped).Stacks)
	assert.NotEmpty(t, (*wrapped).Caller)
}
