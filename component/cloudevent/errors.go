package cloudevent

import (
	"fmt"
	"time"

	"github.com/pubgo/funk/v2/errors"
)

var (
	errReject        = errors.New("cloudevent: reject retry and discard msg")
	errRedeliveryStr = "cloudevent: redelivery message with custom delay duration"
	errForceRetry    = errors.New("cloudevent: force redelivery message with no delay")
)

func Reject(errs ...error) error {
	reason := "reject"
	if len(errs) > 0 {
		reason = errs[0].Error()
	}
	return errors.Wrap(errReject, reason)
}

func isRejectErr(err error) bool {
	return err != nil && errors.Is(err, errReject)
}

type errRedelivery struct {
	delay time.Duration
}

func (err errRedelivery) Error() string {
	return errRedeliveryStr + fmt.Sprintf(":%s", err.delay)
}

func Redelivery(delay time.Duration, errs ...error) error {
	reason := "redelivery"
	if len(errs) > 0 {
		reason = errs[0].Error()
	}
	return errors.Wrap(&errRedelivery{delay: delay}, reason)
}

func isRedeliveryErr(err error) *errRedelivery {
	if err == nil {
		return nil
	}

	var err1 errRedelivery
	if errors.As(err, &err1) {
		return &err1
	}
	return nil
}

type forceRetryError struct{}

func (err forceRetryError) Error() string {
	return errForceRetry.Error()
}

func ForceRetry(errs ...error) error {
	reason := "force_retry"
	if len(errs) > 0 {
		reason = errs[0].Error()
	}
	return errors.Wrap(&forceRetryError{}, reason)
}

func isForceRetry(err error) bool {
	return err != nil && errors.As(err, new(forceRetryError))
}
