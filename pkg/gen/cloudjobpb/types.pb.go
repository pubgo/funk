// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: cloudjobs/types.proto

package cloudjobpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PushEventOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The content type for the data (optional).
	ContentType *string `protobuf:"bytes,1,opt,name=content_type,json=contentType,proto3,oneof" json:"content_type,omitempty"`
	// The metadata passing to pub components
	//
	// metadata property:
	// - key : the key of the message.
	Metadata map[string]string `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The delay duration of the message.
	DelayDur *durationpb.Duration `protobuf:"bytes,3,opt,name=delay_dur,json=delayDur,proto3,oneof" json:"delay_dur,omitempty"`
	// The message id
	MsgId *string `protobuf:"bytes,4,opt,name=msg_id,json=msgId,proto3,oneof" json:"msg_id,omitempty"`
}

func (x *PushEventOptions) Reset() {
	*x = PushEventOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudjobs_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushEventOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushEventOptions) ProtoMessage() {}

func (x *PushEventOptions) ProtoReflect() protoreflect.Message {
	mi := &file_cloudjobs_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushEventOptions.ProtoReflect.Descriptor instead.
func (*PushEventOptions) Descriptor() ([]byte, []int) {
	return file_cloudjobs_types_proto_rawDescGZIP(), []int{0}
}

func (x *PushEventOptions) GetContentType() string {
	if x != nil && x.ContentType != nil {
		return *x.ContentType
	}
	return ""
}

func (x *PushEventOptions) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *PushEventOptions) GetDelayDur() *durationpb.Duration {
	if x != nil {
		return x.DelayDur
	}
	return nil
}

func (x *PushEventOptions) GetMsgId() string {
	if x != nil && x.MsgId != nil {
		return *x.MsgId
	}
	return ""
}

type RegisterJobOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName *string `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3,oneof" json:"job_name,omitempty"`
}

func (x *RegisterJobOptions) Reset() {
	*x = RegisterJobOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cloudjobs_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterJobOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterJobOptions) ProtoMessage() {}

func (x *RegisterJobOptions) ProtoReflect() protoreflect.Message {
	mi := &file_cloudjobs_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterJobOptions.ProtoReflect.Descriptor instead.
func (*RegisterJobOptions) Descriptor() ([]byte, []int) {
	return file_cloudjobs_types_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterJobOptions) GetJobName() string {
	if x != nil && x.JobName != nil {
		return *x.JobName
	}
	return ""
}

var File_cloudjobs_types_proto protoreflect.FileDescriptor

var file_cloudjobs_types_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6c, 0x61, 0x76, 0x61, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x73, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc6, 0x02, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x0c,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x4a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6c, 0x61, 0x76, 0x61, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x3b, 0x0a, 0x09, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x64, 0x75, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x01,
	0x52, 0x08, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x44, 0x75, 0x72, 0x88, 0x01, 0x01, 0x12, 0x1a, 0x0a,
	0x06, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52,
	0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x88, 0x01, 0x01, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x64, 0x65, 0x6c, 0x61,
	0x79, 0x5f, 0x64, 0x75, 0x72, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64,
	0x22, 0x41, 0x0a, 0x12, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4a, 0x6f, 0x62, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1e, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e,
	0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6a, 0x6f, 0x62, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x70, 0x62, 0x3b,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_cloudjobs_types_proto_rawDescOnce sync.Once
	file_cloudjobs_types_proto_rawDescData = file_cloudjobs_types_proto_rawDesc
)

func file_cloudjobs_types_proto_rawDescGZIP() []byte {
	file_cloudjobs_types_proto_rawDescOnce.Do(func() {
		file_cloudjobs_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_cloudjobs_types_proto_rawDescData)
	})
	return file_cloudjobs_types_proto_rawDescData
}

var file_cloudjobs_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_cloudjobs_types_proto_goTypes = []any{
	(*PushEventOptions)(nil),    // 0: lava.cloudjobs.PushEventOptions
	(*RegisterJobOptions)(nil),  // 1: lava.cloudjobs.RegisterJobOptions
	nil,                         // 2: lava.cloudjobs.PushEventOptions.MetadataEntry
	(*durationpb.Duration)(nil), // 3: google.protobuf.Duration
}
var file_cloudjobs_types_proto_depIdxs = []int32{
	2, // 0: lava.cloudjobs.PushEventOptions.metadata:type_name -> lava.cloudjobs.PushEventOptions.MetadataEntry
	3, // 1: lava.cloudjobs.PushEventOptions.delay_dur:type_name -> google.protobuf.Duration
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_cloudjobs_types_proto_init() }
func file_cloudjobs_types_proto_init() {
	if File_cloudjobs_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cloudjobs_types_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*PushEventOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cloudjobs_types_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterJobOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_cloudjobs_types_proto_msgTypes[0].OneofWrappers = []any{}
	file_cloudjobs_types_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cloudjobs_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cloudjobs_types_proto_goTypes,
		DependencyIndexes: file_cloudjobs_types_proto_depIdxs,
		MessageInfos:      file_cloudjobs_types_proto_msgTypes,
	}.Build()
	File_cloudjobs_types_proto = out.File
	file_cloudjobs_types_proto_rawDesc = nil
	file_cloudjobs_types_proto_goTypes = nil
	file_cloudjobs_types_proto_depIdxs = nil
}
