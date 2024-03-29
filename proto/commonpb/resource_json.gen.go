// Code generated by protoc-gen-jsonshim. DO NOT EDIT.
package commonpb

import (
	bytes "bytes"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

// MarshalJSON is a custom marshaler for Identifier
func (this *Identifier) MarshalJSON() ([]byte, error) {
	str, err := ResourceMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for Identifier
func (this *Identifier) UnmarshalJSON(b []byte) error {
	return ResourceUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

var (
	ResourceMarshaler   = &jsonpb.Marshaler{}
	ResourceUnmarshaler = &jsonpb.Unmarshaler{AllowUnknownFields: true}
)
