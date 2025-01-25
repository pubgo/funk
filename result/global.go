package result

import (
	"reflect"
	
	"github.com/pubgo/funk/stack"
)

var errChecks []func(error) error

func RegisterErrCheck(f func(error) error) {
	errChecks = append(errChecks, f)
}

func FindErrCheckList() []*stack.Frame {
	var frames []*stack.Frame
	for _, err := range errChecks {
		frames = append(frames, stack.CallerWithFunc(err))
	}
	return frames
}

func RemoveErrCheck(f func(error) error) {
	var checkFrame = stack.CallerWithFunc(f)
	var index = -1
	for idx, errFunc := range errChecks {
		if reflect.DeepEqual(checkFrame, stack.CallerWithFunc(errFunc)) {
			index = idx
			break
		}
	}

	if index != -1 {
		errChecks = append(errChecks[:index], errChecks[index+1:]...)
	}
}
