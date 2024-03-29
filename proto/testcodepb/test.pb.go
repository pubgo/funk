// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: testcodepb/test.proto

package testcodepb

import (
	_ "github.com/pubgo/funk/proto/errorpb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Code int32

const (
	Code_OK Code = 0
	// NotFound 找不到
	Code_NotFound Code = 100000
	// Unknown 未知
	Code_Unknown Code = 100001
	// db connect error
	Code_DbConn Code = 100003
	// default code
	Code_UnknownCode Code = 100004
	// custom msg, 注释可以做其他的操作
	Code_CustomCode Code = 100005
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0:      "OK",
		100000: "NotFound",
		100001: "Unknown",
		100003: "DbConn",
		100004: "UnknownCode",
		100005: "CustomCode",
	}
	Code_value = map[string]int32{
		"OK":          0,
		"NotFound":    100000,
		"Unknown":     100001,
		"DbConn":      100003,
		"UnknownCode": 100004,
		"CustomCode":  100005,
	}
)

func (x Code) Enum() *Code {
	p := new(Code)
	*p = x
	return p
}

func (x Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Code) Descriptor() protoreflect.EnumDescriptor {
	return file_testcodepb_test_proto_enumTypes[0].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_testcodepb_test_proto_enumTypes[0]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{0}
}

var File_testcodepb_test_proto protoreflect.FileDescriptor

var file_testcodepb_test_proto_rawDesc = []byte{
	0x0a, 0x15, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x2f, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65,
	0x73, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x9c,
	0x01, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12,
	0x16, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xa0, 0x8d, 0x06, 0x1a,
	0x06, 0x92, 0xea, 0x30, 0x02, 0x08, 0x05, 0x12, 0x15, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f,
	0x77, 0x6e, 0x10, 0xa1, 0x8d, 0x06, 0x1a, 0x06, 0x92, 0xea, 0x30, 0x02, 0x08, 0x05, 0x12, 0x14,
	0x0a, 0x06, 0x44, 0x62, 0x43, 0x6f, 0x6e, 0x6e, 0x10, 0xa3, 0x8d, 0x06, 0x1a, 0x06, 0x92, 0xea,
	0x30, 0x02, 0x08, 0x0d, 0x12, 0x11, 0x0a, 0x0b, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x43,
	0x6f, 0x64, 0x65, 0x10, 0xa4, 0x8d, 0x06, 0x12, 0x2a, 0x0a, 0x0a, 0x43, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x43, 0x6f, 0x64, 0x65, 0x10, 0xa5, 0x8d, 0x06, 0x1a, 0x18, 0x92, 0xea, 0x30, 0x14, 0x12,
	0x12, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x20,
	0x6d, 0x73, 0x67, 0x1a, 0x08, 0x8a, 0xea, 0x30, 0x04, 0x08, 0x01, 0x10, 0x0d, 0x42, 0x33, 0x5a,
	0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67,
	0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x73,
	0x74, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x3b, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_testcodepb_test_proto_rawDescOnce sync.Once
	file_testcodepb_test_proto_rawDescData = file_testcodepb_test_proto_rawDesc
)

func file_testcodepb_test_proto_rawDescGZIP() []byte {
	file_testcodepb_test_proto_rawDescOnce.Do(func() {
		file_testcodepb_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_testcodepb_test_proto_rawDescData)
	})
	return file_testcodepb_test_proto_rawDescData
}

var file_testcodepb_test_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_testcodepb_test_proto_goTypes = []interface{}{
	(Code)(0), // 0: demo.test.v1.Code
}
var file_testcodepb_test_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_testcodepb_test_proto_init() }
func file_testcodepb_test_proto_init() {
	if File_testcodepb_test_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_testcodepb_test_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_testcodepb_test_proto_goTypes,
		DependencyIndexes: file_testcodepb_test_proto_depIdxs,
		EnumInfos:         file_testcodepb_test_proto_enumTypes,
	}.Build()
	File_testcodepb_test_proto = out.File
	file_testcodepb_test_proto_rawDesc = nil
	file_testcodepb_test_proto_goTypes = nil
	file_testcodepb_test_proto_depIdxs = nil
}
