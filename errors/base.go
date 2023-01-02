package errors

import (
	"errors"
	"fmt"

	jjson "github.com/goccy/go-json"
)

type Err struct {
	Err    error  `json:"err"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
}

func (e Err) MarshalJSON() ([]byte, error) {
	var data = make(map[string]string, 4)
	data["msg"] = e.Msg
	data["detail"] = e.Detail
	if e.Err != nil {
		data["err"] = e.Err.Error()
		data["err_detail"] = fmt.Sprintf("%#v", e.Err)
	}
	return jjson.Marshal(data)
}

func (e Err) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}

	return errors.New(e.String())
}

func (e Err) String() string {
	return fmt.Sprintf("msg=%q detail=%q error=%q", e.Msg, e.Detail, e.Err)
}

func (e Err) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.String()
}
