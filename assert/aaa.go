package assert

import "github.com/pubgo/funk/v2/features"

const Name = "assert"

var FeatureDebugMode = features.Bool("assert.debug_mode", false, "debug mode, pretty stack and error")
