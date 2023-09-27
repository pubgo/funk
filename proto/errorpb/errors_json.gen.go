// Code generated by protoc-gen-jsonshim. DO NOT EDIT.
package errorpb

import (
	bytes "bytes"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

// MarshalJSON is a custom marshaler for ErrRedirect
func (this *ErrRedirect) MarshalJSON() ([]byte, error) {
	str, err := ErrorsMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for ErrRedirect
func (this *ErrRedirect) UnmarshalJSON(b []byte) error {
	return ErrorsUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for ErrMsg
func (this *ErrMsg) MarshalJSON() ([]byte, error) {
	str, err := ErrorsMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for ErrMsg
func (this *ErrMsg) UnmarshalJSON(b []byte) error {
	return ErrorsUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for ErrCode
func (this *ErrCode) MarshalJSON() ([]byte, error) {
	str, err := ErrorsMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for ErrCode
func (this *ErrCode) UnmarshalJSON(b []byte) error {
	return ErrorsUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for ErrTrace
func (this *ErrTrace) MarshalJSON() ([]byte, error) {
	str, err := ErrorsMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for ErrTrace
func (this *ErrTrace) UnmarshalJSON(b []byte) error {
	return ErrorsUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for Error
func (this *Error) MarshalJSON() ([]byte, error) {
	str, err := ErrorsMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for Error
func (this *Error) UnmarshalJSON(b []byte) error {
	return ErrorsUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

var (
	ErrorsMarshaler   = &jsonpb.Marshaler{}
	ErrorsUnmarshaler = &jsonpb.Unmarshaler{AllowUnknownFields: true}
)
