package errcode

import (
	"fmt"

	"github.com/pubgo/funk/v2/proto/errorpb"
	"google.golang.org/protobuf/proto"
)

var errorCodes = make(map[string]*errorpb.ErrCode)

// GetErrCodes returns all registered error codes.
func GetErrCodes() []*errorpb.ErrCode {
	codeList := make([]*errorpb.ErrCode, 0, len(errorCodes))
	for _, v := range errorCodes {
		codeList = append(codeList, v)
	}
	return codeList
}

// LookupErrCode returns a registered error code by name.
func LookupErrCode(name string) (*errorpb.ErrCode, bool) {
	code, ok := errorCodes[name]
	if !ok {
		return nil, false
	}
	return proto.Clone(code).(*errorpb.ErrCode), true
}

// RegisterErrCode registers an error code and returns an error when the name already exists.
func RegisterErrCode(code *errorpb.ErrCode) error {
	if code == nil {
		return fmt.Errorf("errcode: code is nil")
	}
	if code.Name == "" {
		return fmt.Errorf("errcode: code name is empty")
	}
	if errorCodes[code.Name] != nil {
		return fmt.Errorf("errcode: already registered: name=%q", code.Name)
	}

	errorCodes[code.Name] = proto.Clone(code).(*errorpb.ErrCode)
	return nil
}

// MustRegisterErrCode registers an error code and panics on failure.
func MustRegisterErrCode(code *errorpb.ErrCode) {
	if err := RegisterErrCode(code); err != nil {
		panic(err)
	}
}

// RegisterErrCodes registers an error code and panics when the name already exists.
// Deprecated: use RegisterErrCode or MustRegisterErrCode.
func RegisterErrCodes(code *errorpb.ErrCode) error {
	if err := RegisterErrCode(code); err != nil {
		panic(err)
	}
	return nil
}
