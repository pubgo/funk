package errcode

import (
	"fmt"

	"github.com/pubgo/funk/v2/proto/errorpb"
)

var errorCodes = make(map[string]*errorpb.ErrCode)

func GetErrCodes() []*errorpb.ErrCode {
	var codeList []*errorpb.ErrCode
	for _, v := range errorCodes {
		codeList = append(codeList, v)
	}
	return codeList
}

func RegisterErrCodes(code *errorpb.ErrCode) error {
	if errorCodes[code.Name] != nil {
		panic(fmt.Sprintf("code exists, code=%s", errorCodes[code.Name]))
	}

	errorCodes[code.Name] = code
	return nil
}
