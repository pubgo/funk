package errinter

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Error interface {
	error
	ID() string
	Kind() string
	String() string
	MarshalJSON() ([]byte, error)
}

type ErrUnwrapper interface {
	Unwrap() error
}

type ErrorProto interface {
	error
	Proto() proto.Message
}

type GRPCStatus interface {
	GRPCStatus() *status.Status
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

var (
	_ json.Marshaler = (Tags)(nil)
	_ fmt.Formatter  = (Tags)(nil)
)

type Maps map[string]any

func (t Maps) Tags() Tags {
	tags := make(Tags, 0, len(t))
	for k, v := range t {
		tags = append(tags, Tag{K: k, V: v})
	}
	return tags
}

func (t Maps) ToMapString() map[string]string {
	var data = make(map[string]string, len(t))
	for key, value := range t {
		data[key] = fmt.Sprintf("%v", value)
	}
	return data
}

type Tags []Tag

func (t Tags) Format(f fmt.State, verb rune) {
	tags := make(map[string]any, len(t))
	for i := range t {
		tags[t[i].K] = t[i].V
	}

	data, err := json.Marshal(tags)
	if err != nil {
		fmt.Fprintf(f, "%v", err)
	} else {
		fmt.Fprintln(f, string(data))
	}
}

func (t Tags) ToMap() map[string]string {
	var data = make(map[string]string, len(t))
	for _, tag := range t {
		data[tag.K] = fmt.Sprintf("%v", tag.V)
	}
	return data
}

func (t Tags) MarshalJSON() ([]byte, error) {
	tags := make(map[string]any, len(t))
	for i := range t {
		tags[t[i].K] = t[i].V
	}
	return json.Marshal(tags)
}

type Tag struct {
	K string
	V any
}

func (t Tag) String() string {
	return fmt.Sprintf("%s: %v", t.K, t.V)
}
