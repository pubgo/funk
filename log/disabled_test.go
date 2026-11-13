package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgo/funk/v2/log"
)

func TestSetDisabledLoggers(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		target   string
		want     bool
	}{
		{"empty list disables nothing", nil, "grpc", false},
		{"bare name matches exact", []string{"grpc"}, "grpc", true},
		{"bare name matches subtree", []string{"grpc"}, "grpc.server", true},
		{"bare name matches deep subtree", []string{"grpc"}, "grpc.a.b", true},
		{"bare name does not overmatch", []string{"grpc"}, "grpcx", false},
		{"bare name does not undermatch", []string{"grpc"}, "grp", false},
		{"bare name does not match unrelated", []string{"grpc"}, "sql", false},
		{"explicit single-level star", []string{"grpc.*"}, "grpc.server", true},
		{"explicit single-level star stops at dots", []string{"grpc.*"}, "grpc.a.b", false},
		{"explicit globstar subtree", []string{"grpc.**"}, "grpc.a.b", true},
		{"explicit globstar does not match root", []string{"grpc.**"}, "grpc", false},
		{"globstar suffix at any depth", []string{"**.debug"}, "a.b.debug", true},
		{"single star suffix single depth only", []string{"*.debug"}, "a.b.debug", false},
		{"footgun grpc* matches grpcfoo", []string{"grpc*"}, "grpcfoo", true},
		{"footgun grpc* does not cross dots", []string{"grpc*"}, "grpc.server", false},
		{"single wildcard segment", []string{"*.debug"}, "foo.debug", true},
		{"wildcard does not cross dots", []string{"*.debug"}, "a.b.debug", false},
		{"alternation", []string{"grpc.{server,client}"}, "grpc.client", true},
		{"alternation miss", []string{"grpc.{server,client}"}, "grpc.metrics", false},
		{"blank entries skipped", []string{"", "  "}, "grpc", false},
		{"multiple patterns", []string{"sql", "grpc.*"}, "grpc.server", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.SetDisabledLoggers(tt.patterns)
			t.Cleanup(func() { log.SetDisabledLoggers(nil) })

			assert.Equal(t, tt.want, log.IsDisabled(tt.target))
		})
	}
}

func TestSetDisabledLoggersInvalidPatternPanics(t *testing.T) {
	t.Cleanup(func() { log.SetDisabledLoggers(nil) })

	// Unclosed bracket class is a real glob syntax error; gobwas tolerates
	// other "malformed" inputs (e.g. unclosed braces) as literals.
	assert.Panics(t, func() { log.SetDisabledLoggers([]string{"grpc.[abc"}) })
}

func TestIsDisabledEmptyName(t *testing.T) {
	log.SetDisabledLoggers([]string{"*"})
	t.Cleanup(func() { log.SetDisabledLoggers(nil) })

	assert.False(t, log.IsDisabled(""), "empty name must never be disabled")
}

func TestIsDisabledLoggerField(t *testing.T) {
	log.SetDisabledLoggers([]string{"env"})
	t.Cleanup(func() { log.SetDisabledLoggers(nil) })

	fields := log.Fields{"logger": "env"}
	assert.True(t, log.IsDisabledEvent("env", fields), "field logger=env matches")
	assert.True(t, log.IsDisabledEvent("", fields), "unnamed slog logger matches via field")
	assert.False(t, log.IsDisabledEvent("", log.Fields{"logger": "other"}))
	assert.False(t, log.IsDisabledEvent("", nil), "nil fields are safe")
}
