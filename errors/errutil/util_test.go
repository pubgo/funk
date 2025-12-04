package errutil_test

import (
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/errors/errutil"
)

func TestIsMemoryErr(t *testing.T) {
	// 测试nil错误
	assert.False(t, errutil.IsMemoryErr(nil))

	// 测试内存错误
	memErr := errors.New("runtime error: invalid memory address or nil pointer dereference")
	assert.True(t, errutil.IsMemoryErr(memErr))

	// 测试非内存错误
	normalErr := errors.New("normal error")
	assert.False(t, errutil.IsMemoryErr(normalErr))
}

func TestIsRetryableHTTP(t *testing.T) {
	// 测试GOAWAY错误
	goAwayErr := errors.New("http2: server sent GOAWAY")
	retryType, isRetryable := errutil.IsRetryableHTTP(goAwayErr)
	assert.True(t, isRetryable)
	assert.Equal(t, "goaway", retryType)

	// 测试服务器关闭空闲连接错误
	idleErr := errors.New("http: server closed idle connection")
	retryType, isRetryable = errutil.IsRetryableHTTP(idleErr)
	assert.True(t, isRetryable)
	assert.Equal(t, "server_close_idle_connection", retryType)

	// 测试非重试错误
	normalErr := errors.New("normal http error")
	retryType, isRetryable = errutil.IsRetryableHTTP(normalErr)
	assert.False(t, isRetryable)
	assert.Empty(t, retryType)
}

func TestIsRetryableNetwork(t *testing.T) {
	// 创建一个临时错误
	tempErr := &temporaryError{true}
	retryType, isRetryable := errutil.IsRetryableNetwork(tempErr)
	assert.True(t, isRetryable)
	assert.Equal(t, "temporary", retryType)

	// 创建一个超时错误
	timeoutErr := &timeoutError{true}
	retryType, isRetryable = errutil.IsRetryableNetwork(timeoutErr)
	assert.True(t, isRetryable)
	assert.Equal(t, "timeout", retryType)

	// 测试普通错误
	normalErr := errors.New("normal network error")
	retryType, isRetryable = errutil.IsRetryableNetwork(normalErr)
	assert.False(t, isRetryable)
	assert.Empty(t, retryType)
}

func TestIsTemporary(t *testing.T) {
	// 测试临时错误
	tempErr := &temporaryError{true}
	assert.True(t, errutil.IsTemporary(tempErr))

	// 测试非临时错误
	nonTempErr := &temporaryError{false}
	assert.False(t, errutil.IsTemporary(nonTempErr))

	// 测试普通错误
	normalErr := errors.New("normal error")
	assert.False(t, errutil.IsTemporary(normalErr))
}

func TestIsTimeout(t *testing.T) {
	// 测试超时错误
	timeoutErr := &timeoutError{true}
	assert.True(t, errutil.IsTimeout(timeoutErr))

	// 测试非超时错误
	nonTimeoutErr := &timeoutError{false}
	assert.False(t, errutil.IsTimeout(nonTimeoutErr))

	// 测试普通错误
	normalErr := errors.New("normal error")
	assert.False(t, errutil.IsTimeout(normalErr))
}

func TestIsTemporaryConnection(t *testing.T) {
	// 在某些系统上，EAGAIN和EWOULDBLOCK可能是相同的值，所以我们需要特殊处理
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ECONNRESET", syscall.ECONNRESET, true},
		{"ECONNABORTED", syscall.ECONNABORTED, true},
		{"ENOTCONN", syscall.ENOTCONN, true},
		{"ETIMEDOUT", syscall.ETIMEDOUT, true},
		{"EINTR", syscall.EINTR, true},
		{"normal error", errors.New("normal error"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, isRetryable := errutil.IsTemporaryConnection(tc.err)
			assert.Equal(t, tc.expected, isRetryable)
		})
	}

	// 单独测试EAGAIN和EWOULDBLOCK
	_, isRetryable := errutil.IsTemporaryConnection(syscall.EAGAIN)
	assert.True(t, isRetryable)

	_, isRetryable = errutil.IsTemporaryConnection(syscall.EWOULDBLOCK)
	assert.True(t, isRetryable)
}

// 辅助类型用于测试
type temporaryError struct {
	isTemporary bool
}

func (e *temporaryError) Error() string {
	return "temporary error"
}

func (e *temporaryError) Temporary() bool {
	return e.isTemporary
}

type timeoutError struct {
	isTimeout bool
}

func (e *timeoutError) Error() string {
	return "timeout error"
}

func (e *timeoutError) Timeout() bool {
	return e.isTimeout
}
