package recovery

import (
	"os"
	"testing"
)

// SetExitFn replaces the process exit hook. It is intended for tests.
func SetExitFn(fn func(code int)) {
	if fn == nil {
		exitFn = os.Exit
		return
	}
	exitFn = fn
}

// SetTestingFatalFn replaces the fatal hook used by Testing. It is intended for tests.
func SetTestingFatalFn(fn func(t *testing.T, err error)) {
	if fn == nil {
		testingFatalFn = func(t *testing.T, err error) {
			t.Fatal(err)
		}
		return
	}
	testingFatalFn = fn
}
