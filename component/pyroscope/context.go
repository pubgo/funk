package pyroscope

import (
	"context"

	"github.com/grafana/pyroscope-go"
)

type LabelSet = pyroscope.LabelSet

var Labels = pyroscope.Labels

func TagWrapper(ctx context.Context, labels LabelSet, cb func(context.Context)) {
	pyroscope.TagWrapper(ctx, labels, cb)
}
