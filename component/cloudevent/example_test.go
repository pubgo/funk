package cloudevent_test

import (
	"fmt"

	"github.com/pubgo/funk/v2/component/cloudevent"
	cloudeventpb "github.com/pubgo/funk/v2/proto/cloudevent"
	"github.com/samber/lo"
)

func ExampleWithPushOpt() {
	opt := cloudevent.WithPushOpt(func(o *cloudeventpb.PushEventOptions) {
		o.ContentType = lo.ToPtr("application/protobuf")
	})
	fmt.Println(opt.GetContentType())

	// Output: application/protobuf
}

func ExampleProtoRegisterOpts() {
	opt := cloudevent.ProtoRegisterOpts(&cloudeventpb.RegisterJobOptions{
		JobName: lo.ToPtr("demo"),
	})
	ro := &cloudevent.RegisterJobOptions{Opts: new(cloudeventpb.RegisterJobOptions)}
	opt(ro)
	fmt.Println(lo.FromPtr(ro.Opts.JobName))

	// Output: demo
}
