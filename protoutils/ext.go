package protoutils

import (
	"github.com/pubgo/funk/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func HasExtension(option protoreflect.ProtoMessage, xt protoreflect.ExtensionType) bool {
	if option == nil || xt == nil {
		return false
	}

	return proto.HasExtension(option, xt)
}

func GetExtension[Option any](option protoreflect.ProtoMessage, xt protoreflect.ExtensionType) *Option {
	if option == nil || xt == nil {
		return nil
	}

	opt, ok := proto.GetExtension(option, xt).(*Option)
	if !ok || opt == nil {
		return nil
	}

	return opt
}

func EachServiceMethod(srv protoreflect.ServiceDescriptor, fn func(mth protoreflect.MethodDescriptor)) {
	mthLen := srv.Methods().Len()
	for i := 0; i < mthLen; i++ {
		fn(srv.Methods().Get(i))
	}
}

func EachService(fn func(desc protoreflect.FileDescriptor, srv protoreflect.ServiceDescriptor)) {
	protoregistry.GlobalFiles.RangeFiles(func(desc protoreflect.FileDescriptor) bool {
		srvLen := desc.Services().Len()
		for i := 0; i < srvLen; i++ {
			fn(desc, desc.Services().Get(i))
		}
		return true
	})
}

func GetService(name string) protoreflect.ServiceDescriptor {
	service, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil {
		log.Err(err).Str("name", name).Msg("failed to find protobuf service")
		return nil
	}

	srv, ok := service.(protoreflect.ServiceDescriptor)
	if ok {
		return srv
	}
	return nil
}
