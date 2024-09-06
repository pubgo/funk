package errors

import (
	"fmt"

	"github.com/pubgo/funk/proto/errorpb"
)

var codes = make(map[string]*errorpb.ErrCode)

func GetErrCodes() []*errorpb.ErrCode {
	var codeList []*errorpb.ErrCode
	for _, v := range codes {
		codeList = append(codeList, v)
	}
	return codeList
}

func RegisterErrCodes(code *errorpb.ErrCode) error {
	if codes[code.Name] != nil {
		panic(fmt.Sprintf("code exists, code=%s", codes[code.Name]))
	}

	codes[code.Name] = code
	return nil
}
