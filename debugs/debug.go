package debugs

import "github.com/pubgo/funk/v2/features"

var Enabled = features.Bool("debug.enabled", false, "feature: enable debug mode")

func SetEnabled() {
	_ = Enabled.Set("true")
}

func SetDisabled() {
	_ = Enabled.Set("false")
}
