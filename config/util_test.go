package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveYAMLComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple comment removal",
			input:    "key: value # this is a comment\nother: value",
			expected: "key: value\nother: value",
		},
		{
			name:     "comment at beginning of line",
			input:    "# this is a comment\nkey: value",
			expected: "key: value",
		},
		{
			name:     "quoted string with hash",
			input:    "key: \"value # not a comment\"\nother: value",
			expected: "key: \"value # not a comment\"\nother: value",
		},
		{
			name:     "hash inside quoted string",
			input:    "key: \"#comment here\" # real comment\nother: value",
			expected: "key: \"#comment here\"\nother: value",
		},
		{
			name:     "single quoted string with hash",
			input:    "key: '#not a comment'\nother: value",
			expected: "key: '#not a comment'\nother: value",
		},
		{
			name:     "mixed quotes",
			input:    "key: \"double quote # comment\"\nother: 'single quote # comment'\nthird: value # actual comment",
			expected: "key: \"double quote # comment\"\nother: 'single quote # comment'\nthird: value",
		},
		{
			name:     "escaped quotes not handled specially",
			input:    "key: \"value with \\\"quotes\\\" # not a comment\"\nother: value",
			expected: "key: \"value with \\\"quotes\\\" # not a comment\"\nother: value",
		},
		{
			name:     "empty lines preserved",
			input:    "key: value\n\nother: value",
			expected: "key: value\n\nother: value",
		},
		{
			name:     "multiple comments",
			input:    "key: value # comment1\nother: value # comment2",
			expected: "key: value\nother: value",
		},
		{
			name:     "whitespace before comment",
			input:    "key: value   # comment with spaces\nother: value",
			expected: "key: value\nother: value",
		},
		{
			name:     "only comment line",
			input:    "key: value\n  # this is a comment\nother: value",
			expected: "key: value\nother: value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeYAMLComments([]byte(tt.input))
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestRemoveYAMLCommentsFromLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple comment",
			input:    "key: value # comment",
			expected: "key: value",
		},
		{
			name:     "quoted hash",
			input:    "key: \"value # not comment\"",
			expected: "key: \"value # not comment\"",
		},
		{
			name:     "comment after quoted hash",
			input:    "key: \"value # not comment\" # actual comment",
			expected: "key: \"value # not comment\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeYAMLCommentsFromLine([]byte(tt.input))
			assert.Equal(t, tt.expected, string(result))
		})
	}
}
