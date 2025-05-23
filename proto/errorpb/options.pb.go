// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.3
// source: errorpb/options.proto

package errorpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gen         bool   `protobuf:"varint,1,opt,name=gen,proto3" json:"gen,omitempty"`
	DefaultCode Code   `protobuf:"varint,2,opt,name=default_code,json=defaultCode,proto3,enum=errors.Code" json:"default_code,omitempty"`
	Name        string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Options) Reset() {
	*x = Options{}
	mi := &file_errorpb_options_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Options) ProtoMessage() {}

func (x *Options) ProtoReflect() protoreflect.Message {
	mi := &file_errorpb_options_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Options.ProtoReflect.Descriptor instead.
func (*Options) Descriptor() ([]byte, []int) {
	return file_errorpb_options_proto_rawDescGZIP(), []int{0}
}

func (x *Options) GetGen() bool {
	if x != nil {
		return x.Gen
	}
	return false
}

func (x *Options) GetDefaultCode() Code {
	if x != nil {
		return x.DefaultCode
	}
	return Code_OK
}

func (x *Options) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Fields struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code Code   `protobuf:"varint,1,opt,name=code,proto3,enum=errors.Code" json:"code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Fields) Reset() {
	*x = Fields{}
	mi := &file_errorpb_options_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Fields) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Fields) ProtoMessage() {}

func (x *Fields) ProtoReflect() protoreflect.Message {
	mi := &file_errorpb_options_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Fields.ProtoReflect.Descriptor instead.
func (*Fields) Descriptor() ([]byte, []int) {
	return file_errorpb_options_proto_rawDescGZIP(), []int{1}
}

func (x *Fields) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_OK
}

func (x *Fields) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Fields) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var file_errorpb_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*Fields)(nil),
		Field:         100002,
		Name:          "errors.field",
		Tag:           "bytes,100002,opt,name=field",
		Filename:      "errorpb/options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         100001,
		Name:          "errors.opts",
		Tag:           "bytes,100001,opt,name=opts",
		Filename:      "errorpb/options.proto",
	},
}

// Extension fields to descriptorpb.EnumValueOptions.
var (
	// optional errors.Fields field = 100002;
	E_Field = &file_errorpb_options_proto_extTypes[0]
)

// Extension fields to descriptorpb.EnumOptions.
var (
	// optional errors.Options opts = 100001;
	E_Opts = &file_errorpb_options_proto_extTypes[1]
)

var File_errorpb_options_proto protoreflect.FileDescriptor

var file_errorpb_options_proto_rawDesc = []byte{
	0x0a, 0x15, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x1a,
	0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x12, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x60, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x67, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x67,
	0x65, 0x6e, 0x12, 0x2f, 0x0a, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x0b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x50, 0x0a, 0x06, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x12, 0x20, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0c, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x49, 0x0a, 0x05, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x12, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x8d, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x52, 0x05, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x3a, 0x43, 0x0a, 0x04, 0x6f, 0x70, 0x74, 0x73, 0x12, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6e, 0x75, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa1, 0x8d, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x04, 0x6f, 0x70, 0x74, 0x73, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x2f, 0x66, 0x75,
	0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62,
	0x3b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errorpb_options_proto_rawDescOnce sync.Once
	file_errorpb_options_proto_rawDescData = file_errorpb_options_proto_rawDesc
)

func file_errorpb_options_proto_rawDescGZIP() []byte {
	file_errorpb_options_proto_rawDescOnce.Do(func() {
		file_errorpb_options_proto_rawDescData = protoimpl.X.CompressGZIP(file_errorpb_options_proto_rawDescData)
	})
	return file_errorpb_options_proto_rawDescData
}

var file_errorpb_options_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_errorpb_options_proto_goTypes = []any{
	(*Options)(nil),                       // 0: errors.Options
	(*Fields)(nil),                        // 1: errors.Fields
	(Code)(0),                             // 2: errors.Code
	(*descriptorpb.EnumValueOptions)(nil), // 3: google.protobuf.EnumValueOptions
	(*descriptorpb.EnumOptions)(nil),      // 4: google.protobuf.EnumOptions
}
var file_errorpb_options_proto_depIdxs = []int32{
	2, // 0: errors.Options.default_code:type_name -> errors.Code
	2, // 1: errors.Fields.code:type_name -> errors.Code
	3, // 2: errors.field:extendee -> google.protobuf.EnumValueOptions
	4, // 3: errors.opts:extendee -> google.protobuf.EnumOptions
	1, // 4: errors.field:type_name -> errors.Fields
	0, // 5: errors.opts:type_name -> errors.Options
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	4, // [4:6] is the sub-list for extension type_name
	2, // [2:4] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_errorpb_options_proto_init() }
func file_errorpb_options_proto_init() {
	if File_errorpb_options_proto != nil {
		return
	}
	file_errorpb_code_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_errorpb_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_errorpb_options_proto_goTypes,
		DependencyIndexes: file_errorpb_options_proto_depIdxs,
		MessageInfos:      file_errorpb_options_proto_msgTypes,
		ExtensionInfos:    file_errorpb_options_proto_extTypes,
	}.Build()
	File_errorpb_options_proto = out.File
	file_errorpb_options_proto_rawDesc = nil
	file_errorpb_options_proto_goTypes = nil
	file_errorpb_options_proto_depIdxs = nil
}
