package debugs

import "github.com/pubgo/funk/v2/features"

var Enabled = features.Bool("debug.enabled", false, "feature: enable debug mode")

// SetEnabled enables debug mode
func SetEnabled() {
	_ = Enabled.Set("true")
}

// SetDisabled disables debug mode
func SetDisabled() {
	_ = Enabled.Set("false")
}

// IsDebug returns true if debug mode is enabled
func IsDebug() bool { return Enabled.Value() }
