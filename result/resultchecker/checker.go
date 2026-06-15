package resultchecker

import (
	"context"
	"slices"
	"sync"

	"github.com/pubgo/funk/v2/stack"
)

var (
	errChecksMu sync.RWMutex
	errChecks   []ErrChecker
)

func RegisterErrCheck(f ErrChecker) bool {
	if f == nil {
		return false
	}

	checkFrame := stack.CallerWithFunc(f).String()

	errChecksMu.Lock()
	defer errChecksMu.Unlock()

	for _, errFunc := range errChecks {
		if checkFrame == stack.CallerWithFunc(errFunc).String() {
			return false
		}
	}

	errChecks = append(errChecks, f)
	return true
}

func GetErrChecks() []ErrChecker {
	errChecksMu.RLock()
	defer errChecksMu.RUnlock()

	return slices.Clone(errChecks)
}

func GetErrCheckStacks() []*stack.Frame {
	errChecksMu.RLock()
	defer errChecksMu.RUnlock()

	frames := make([]*stack.Frame, 0, len(errChecks))
	for _, err := range errChecks {
		frames = append(frames, stack.CallerWithFunc(err))
	}
	return frames
}

func RemoveErrCheck(f func(context.Context, error) error) {
	checkFrame := stack.CallerWithFunc(f).String()

	errChecksMu.Lock()
	defer errChecksMu.Unlock()

	index := -1
	for idx, errFunc := range errChecks {
		if checkFrame == stack.CallerWithFunc(errFunc).String() {
			index = idx
			break
		}
	}

	if index != -1 {
		errChecks = append(errChecks[:index], errChecks[index+1:]...)
	}
}
