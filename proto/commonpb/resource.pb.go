// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: commonpb/resource.proto

package commonpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VISIABLE_MODE int32

const (
	// 默认Namespace可见
	VISIABLE_MODE_NAMESPACE VISIABLE_MODE = 0
	// 域内可见
	VISIABLE_MODE_DOMAIN VISIABLE_MODE = 1
	// 全局可见
	VISIABLE_MODE_GLOBAL VISIABLE_MODE = 2
)

// Enum value maps for VISIABLE_MODE.
var (
	VISIABLE_MODE_name = map[int32]string{
		0: "NAMESPACE",
		1: "DOMAIN",
		2: "GLOBAL",
	}
	VISIABLE_MODE_value = map[string]int32{
		"NAMESPACE": 0,
		"DOMAIN":    1,
		"GLOBAL":    2,
	}
)

func (x VISIABLE_MODE) Enum() *VISIABLE_MODE {
	p := new(VISIABLE_MODE)
	*p = x
	return p
}

func (x VISIABLE_MODE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VISIABLE_MODE) Descriptor() protoreflect.EnumDescriptor {
	return file_commonpb_resource_proto_enumTypes[0].Descriptor()
}

func (VISIABLE_MODE) Type() protoreflect.EnumType {
	return &file_commonpb_resource_proto_enumTypes[0]
}

func (x VISIABLE_MODE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VISIABLE_MODE.Descriptor instead.
func (VISIABLE_MODE) EnumDescriptor() ([]byte, []int) {
	return file_commonpb_resource_proto_rawDescGZIP(), []int{0}
}

// A number of application services require a mechanism to
// encode the identity of a particular resource across applications.
// The resource identifiers consist of the application ID,
// an application-defined resource type,
// and an application-defined ID for that resource
type Identifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// An application identifier that will be used among
	// other infrastructure services to identify the application
	ApplicationName string `protobuf:"bytes,1,opt,name=application_name,json=applicationName,proto3" json:"application_name,omitempty"`
	// An application specific type name of a resource
	ResourceType string `protobuf:"bytes,2,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	// An application specific resource identity of a resource
	ResourceId string `protobuf:"bytes,3,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
}

func (x *Identifier) Reset() {
	*x = Identifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commonpb_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Identifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Identifier) ProtoMessage() {}

func (x *Identifier) ProtoReflect() protoreflect.Message {
	mi := &file_commonpb_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Identifier.ProtoReflect.Descriptor instead.
func (*Identifier) Descriptor() ([]byte, []int) {
	return file_commonpb_resource_proto_rawDescGZIP(), []int{0}
}

func (x *Identifier) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *Identifier) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *Identifier) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

var File_commonpb_resource_proto protoreflect.FileDescriptor

var file_commonpb_resource_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x22, 0x7d, 0x0a, 0x0a, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69,
	0x65, 0x72, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x64, 0x2a, 0x36, 0x0a, 0x0d, 0x56, 0x49, 0x53, 0x49, 0x41, 0x42, 0x4c, 0x45, 0x5f,
	0x4d, 0x4f, 0x44, 0x45, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43,
	0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x10, 0x01, 0x12,
	0x0a, 0x0a, 0x06, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x10, 0x02, 0x42, 0x2f, 0x5a, 0x2d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x2f,
	0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x70, 0x62, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_commonpb_resource_proto_rawDescOnce sync.Once
	file_commonpb_resource_proto_rawDescData = file_commonpb_resource_proto_rawDesc
)

func file_commonpb_resource_proto_rawDescGZIP() []byte {
	file_commonpb_resource_proto_rawDescOnce.Do(func() {
		file_commonpb_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_commonpb_resource_proto_rawDescData)
	})
	return file_commonpb_resource_proto_rawDescData
}

var file_commonpb_resource_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_commonpb_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_commonpb_resource_proto_goTypes = []interface{}{
	(VISIABLE_MODE)(0), // 0: common.v1.VISIABLE_MODE
	(*Identifier)(nil), // 1: common.v1.Identifier
}
var file_commonpb_resource_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_commonpb_resource_proto_init() }
func file_commonpb_resource_proto_init() {
	if File_commonpb_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_commonpb_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Identifier); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_commonpb_resource_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commonpb_resource_proto_goTypes,
		DependencyIndexes: file_commonpb_resource_proto_depIdxs,
		EnumInfos:         file_commonpb_resource_proto_enumTypes,
		MessageInfos:      file_commonpb_resource_proto_msgTypes,
	}.Build()
	File_commonpb_resource_proto = out.File
	file_commonpb_resource_proto_rawDesc = nil
	file_commonpb_resource_proto_goTypes = nil
	file_commonpb_resource_proto_depIdxs = nil
}