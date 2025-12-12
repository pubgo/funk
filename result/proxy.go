package result

import (
	"fmt"

	"github.com/samber/lo"
)

var _ ErrSetter = new(ProxyErr)

type ProxyErr struct {
	err *error
}

func (e ProxyErr) IsOK() bool {
	return lo.FromPtr(e.err) == nil
}

func (e ProxyErr) IsErr() bool {
	return lo.FromPtr(e.err) != nil
}

func (e ProxyErr) GetErr() error {
	return lo.FromPtr(e.err)
}

func (e ProxyErr) Err() error {
	return lo.FromPtr(e.err)
}

func (e ProxyErr) String() string {
	if e.IsOK() {
		return "OK"
	}

	return fmt.Sprintf("Error(%v)", lo.FromPtr(e.err))
}

func (e ProxyErr) setErrorInner() {
}
