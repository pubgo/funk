package result

import (
	"fmt"

	"github.com/samber/lo"
)

var _ ErrSetter = new(ErrProxy)

type ErrProxy struct {
	err *error
}

func (e ErrProxy) IsOK() bool {
	return lo.FromPtr(e.err) == nil
}

func (e ErrProxy) IsErr() bool {
	return lo.FromPtr(e.err) != nil
}

func (e ErrProxy) GetErr() error {
	return lo.FromPtr(e.err)
}

func (e ErrProxy) String() string {
	if e.IsOK() {
		return "Ok"
	}

	return fmt.Sprintf("Error(%v)", lo.FromPtr(e.err))
}

func (e ErrProxy) setError(err error) {
	if err == nil {
		return
	}

	e.err = lo.ToPtr(err)
}
