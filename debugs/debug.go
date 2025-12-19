package debugs

import "github.com/pubgo/funk/v2/features"

var Enabled = features.Bool("debug.enabled", false, "debug mode feature")
