package errors

import (
	"fmt"
	"google.golang.org/grpc/status"

	json "github.com/goccy/go-json"
)

var _ json.Marshaler = (Tags)(nil)
var _ fmt.Formatter = (Tags)(nil)

type Maps map[string]any

func (t Maps) Tags() Tags {
	var tags = make(Tags, 0, len(t))
	for k, v := range t {
		tags = append(tags, Tag{K: k, V: v})
	}
	return tags
}

type Tags []Tag

func (t Tags) Format(f fmt.State, verb rune) {
	var tags = make(map[string]any, len(t))
	for i := range t {
		tags[t[i].K] = t[i].V
	}

	var data, err = json.Marshal(tags)
	if err != nil {
		fmt.Fprintf(f, "%v", err)
	} else {
		fmt.Fprintln(f, string(data))
	}
}

func (t Tags) MarshalJSON() ([]byte, error) {
	var tags = make(map[string]any, len(t))
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

type ErrUnwrap interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type Error interface {
	Kind() string
	Error() string
	String() string
	MarshalJSON() ([]byte, error)
}

type GRPCStatus interface {
	GRPCStatus() *status.Status
}
