package error1

import "github.com/pubgo/funk/proto/errorpb"

type ErrWrap struct {
	err  error
	wrap *errorpb.ErrWrap
}

func (e *ErrWrap) Error() string {
	return e.err.Error()
}

func (e *ErrWrap) String() string {
	return e.wrap.String()
}

func (e *ErrWrap) MarshalJSON() ([]byte, error) {
	return e.wrap.MarshalJSON()
}

func (e *ErrWrap) Unwrap() error {
	return e.err
}

func (e *ErrWrap) Proto() *errorpb.ErrWrap {
	return e.wrap
}
