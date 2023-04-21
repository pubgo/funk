package errors

import (
	"fmt"
)

var _ Error = (*ErrBase)(nil)
var _ fmt.Formatter = (*ErrBase)(nil)

type ErrBase struct {
	err error
}

func (t *ErrBase) Format(f fmt.State, verb rune) {
	panic("implement me")
}

func (t *ErrBase) Kind() string {
	return "err_base"
}

func (t *ErrBase) String() string {
	panic("implement me")
}

func (t *ErrBase) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (t *ErrBase) getData() map[string]any {
	var data = make(map[string]any)
	var mm = errJsonify(t.err)
	if mm != nil {
		for k, v := range mm {
			data[k] = v
		}
	}

	return data
}

func (t *ErrBase) Unwrap() error { return t.err }
func (t *ErrBase) Error() string { return t.err.Error() }
