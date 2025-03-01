package aherrcheck

import (
	"context"
	"reflect"

	"github.com/pubgo/funk/stack"
)

var errChecks []func(context.Context, error) error

func RegisterErrCheck(f func(context.Context, error) error) {
	var checkFrame = stack.CallerWithFunc(f)
	for _, errFunc := range errChecks {
		if reflect.DeepEqual(checkFrame, stack.CallerWithFunc(errFunc)) {
			return
		}
	}

	errChecks = append(errChecks, f)
}

func GetErrChecks() []func(context.Context, error) error { return errChecks }

func GetErrCheckFrames() []*stack.Frame {
	var frames []*stack.Frame
	for _, err := range errChecks {
		frames = append(frames, stack.CallerWithFunc(err))
	}
	return frames
}

func RemoveErrCheck(f func(context.Context, error) error) {
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
