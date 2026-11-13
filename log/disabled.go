package log

import (
	"strings"
	"sync/atomic"

	"github.com/gobwas/glob"

	"github.com/pubgo/funk/v2/assert"
	"github.com/pubgo/funk/v2/log/logfields"
)

// disabledEntry is one compiled disable-list pattern.
type disabledEntry struct {
	raw   string
	match func(target string) bool
}

// disabledList holds the compiled entries, replaced atomically on Set.
var disabledList atomic.Value // []disabledEntry

// SetDisabledLoggers replaces the disabled-logger list with glob patterns.
//
// Logger names are dot-separated hierarchies and patterns compile with '.'
// as the separator, so `*` stays within one component while `**` crosses
// components:
//
//	grpc                  exact "grpc" plus its whole subtree at any depth
//	                      (grpc.server, grpc.a.b, ...) via implicit .**
//	grpc.*                the subtree one level deep only, not "grpc" itself
//	grpc.**               the subtree at any depth, not "grpc" itself
//	**.debug              names ending in ".debug" at any depth
//	grpc.{server,client}  alternation
//
// Patterns without glob magic characters (* ? [ {) get the any-depth subtree
// sugar above; anything else is compiled literally. An invalid pattern panics
// here, at configuration time - never on the log hot path.
func SetDisabledLoggers(patterns []string) {
	entries := make([]disabledEntry, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if strings.ContainsAny(p, "*?[{") {
			g := assert.Must1(glob.Compile(p, '.'))
			entries = append(entries, disabledEntry{raw: p, match: g.Match})
			continue
		}

		// Magic-free entry: exact name plus its subtree at any depth.
		sub := assert.Must1(glob.Compile(p+".**", '.'))
		entries = append(entries, disabledEntry{raw: p, match: func(target string) bool {
			return target == p || sub.Match(target)
		}})
	}

	disabledList.Store(entries)
}

// IsDisabled reports whether name matches the disabled-logger list.
func IsDisabled(name string) bool {
	if name == "" {
		return false
	}

	entries, _ := disabledList.Load().([]disabledEntry)
	for _, e := range entries {
		if e.match(name) {
			return true
		}
	}
	return false
}

// IsDisabledEvent reports whether a log event is disabled by its component
// name or by the logfields.Logger attribute (set by loggers that go through
// slog without a component name, e.g. the env package).
//
// Custom SetEnableChecker implementations should call this first so the
// disabled list keeps working when a checker is installed.
func IsDisabledEvent(name string, fields Fields) bool {
	if IsDisabled(name) {
		return true
	}

	if v, ok := fields[logfields.Logger].(string); ok {
		return IsDisabled(v)
	}

	return false
}
