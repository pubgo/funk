package errors

import (
	"fmt"
)

type ErrorID interface {
	error
	ID() string
}

type Error interface {
	ErrorID
	String() string
	MarshalJSON() ([]byte, error)
}

type ErrUnwrapper interface {
	Unwrap() error
}

type ErrIs interface {
	Is(error) bool
}

type ErrAs interface {
	As(any) bool
}

type Tags map[string]any

func cloneTags(tags Tags) Tags {
	if len(tags) == 0 {
		return nil
	}

	out := make(Tags, len(tags))
	for key, value := range tags {
		out[key] = value
	}
	return out
}

func mergeTags(tags ...Tags) Tags {
	tagList := make(Tags)
	for _, t := range tags {
		for key, value := range t {
			tagList[key] = value
		}
	}
	return tagList
}

func (t Tags) Clone() Tags {
	return cloneTags(t)
}

func (t Tags) Merge(tags Tags) Tags {
	tagList := make(Tags, len(t)+len(tags))

	for key, value := range tags {
		tagList[key] = value
	}

	for key, value := range t {
		tagList[key] = value
	}
	return tagList
}

func (t Tags) ToMapString() map[string]string {
	data := make(map[string]string, len(t))
	for key, value := range t {
		data[key] = fmt.Sprintf("%v", value)
	}
	return data
}
