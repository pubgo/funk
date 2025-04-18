package cloudevent

import (
	"fmt"
	"time"

	"github.com/pubgo/funk/errors"
)

var errReject = errors.New("cloudevent: reject retry and discard msg")
var errRedeliveryStr = "cloudevent: redelivery message with custom delay duration"

func Reject(errs ...error) error {
	var reason = "reject"
	if len(errs) > 0 {
		reason = errs[0].Error()
	}
	return errors.Wrap(errReject, reason)
}

func isRejectErr(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, errReject)
}

type errRedelivery struct {
	delay time.Duration
}

func (err errRedelivery) Error() string {
	return errRedeliveryStr + fmt.Sprintf(":%s", err.delay)
}

func Redelivery(delay time.Duration, errs ...error) error {
	var reason = "redelivery"
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

// TODO force retry
