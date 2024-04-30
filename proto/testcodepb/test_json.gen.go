// Code generated by protoc-gen-jsonshim. DO NOT EDIT.
package testcodepb

import (
	bytes "bytes"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

// MarshalJSON is a custom marshaler for User
func (this *User) MarshalJSON() ([]byte, error) {
	str, err := TestMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for User
func (this *User) UnmarshalJSON(b []byte) error {
	return TestUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for User_Role
func (this *User_Role) MarshalJSON() ([]byte, error) {
	str, err := TestMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for User_Role
func (this *User_Role) UnmarshalJSON(b []byte) error {
	return TestUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

var (
	TestMarshaler   = &jsonpb.Marshaler{}
	TestUnmarshaler = &jsonpb.Unmarshaler{AllowUnknownFields: true}
)
