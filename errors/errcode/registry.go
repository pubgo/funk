package errcode

import (
	"fmt"

	"github.com/pubgo/funk/v2/proto/errorpb"
)

var errorCodes = make(map[string]*errorpb.ErrCode)

// GetErrCodes 获取所有错误码
func GetErrCodes() []*errorpb.ErrCode {
	var codeList = make([]*errorpb.ErrCode, 0, len(errorCodes))
	for _, v := range errorCodes {
		codeList = append(codeList, v)
	}
	return codeList
}

// RegisterErrCodes 注册错误码
func RegisterErrCodes(code *errorpb.ErrCode) error {
	if errorCodes[code.Name] != nil {
		panic(fmt.Sprintf("code exists, code=%s", errorCodes[code.Name]))
	}

	errorCodes[code.Name] = code
	return nil
}
