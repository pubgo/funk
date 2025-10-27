package resultchecker

import (
	"context"

	"github.com/pubgo/funk/v2/stack"
)

var errChecks []ErrChecker

func RegisterErrCheck(f ErrChecker) bool {
	if f == nil {
		return false
	}

	checkFrame := stack.CallerWithFunc(f).String()
	for _, errFunc := range errChecks {
		if checkFrame == stack.CallerWithFunc(errFunc).String() {
			return false
		}
	}

	errChecks = append(errChecks, f)
	return true
}

func GetErrChecks() []ErrChecker { return errChecks }

func GetErrCheckStacks() []*stack.Frame {
	var frames []*stack.Frame
	for _, err := range errChecks {
		frames = append(frames, stack.CallerWithFunc(err))
	}
	return frames
}

func RemoveErrCheck(f func(context.Context, error) error) {
	checkFrame := stack.CallerWithFunc(f).String()
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
