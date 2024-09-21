// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: errorpb/code.proto

package errorpb

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

type Code int32

const (
	Code_OK                 Code = 0
	Code_Canceled           Code = 1
	Code_Unknown            Code = 2
	Code_InvalidArgument    Code = 3
	Code_DeadlineExceeded   Code = 4
	Code_NotFound           Code = 5
	Code_AlreadyExists      Code = 6
	Code_PermissionDenied   Code = 7
	Code_ResourceExhausted  Code = 8
	Code_FailedPrecondition Code = 9
	Code_Aborted            Code = 10
	Code_OutOfRange         Code = 11
	Code_Unimplemented      Code = 12
	Code_Internal           Code = 13
	Code_Unavailable        Code = 14
	Code_DataLoss           Code = 15
	Code_Unauthenticated    Code = 16
	Code_TooManyRequests    Code = 17
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0:  "OK",
		1:  "Canceled",
		2:  "Unknown",
		3:  "InvalidArgument",
		4:  "DeadlineExceeded",
		5:  "NotFound",
		6:  "AlreadyExists",
		7:  "PermissionDenied",
		8:  "ResourceExhausted",
		9:  "FailedPrecondition",
		10: "Aborted",
		11: "OutOfRange",
		12: "Unimplemented",
		13: "Internal",
		14: "Unavailable",
		15: "DataLoss",
		16: "Unauthenticated",
		17: "TooManyRequests",
	}
	Code_value = map[string]int32{
		"OK":                 0,
		"Canceled":           1,
		"Unknown":            2,
		"InvalidArgument":    3,
		"DeadlineExceeded":   4,
		"NotFound":           5,
		"AlreadyExists":      6,
		"PermissionDenied":   7,
		"ResourceExhausted":  8,
		"FailedPrecondition": 9,
		"Aborted":            10,
		"OutOfRange":         11,
		"Unimplemented":      12,
		"Internal":           13,
		"Unavailable":        14,
		"DataLoss":           15,
		"Unauthenticated":    16,
		"TooManyRequests":    17,
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
	return file_errorpb_code_proto_enumTypes[0].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_errorpb_code_proto_enumTypes[0]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_errorpb_code_proto_rawDescGZIP(), []int{0}
}

var File_errorpb_code_proto protoreflect.FileDescriptor

var file_errorpb_code_proto_rawDesc = []byte{
	0x0a, 0x12, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2a, 0xc1, 0x02, 0x0a,
	0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x0c, 0x0a,
	0x08, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x65, 0x64, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x55,
	0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x6e, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x41, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x10, 0x03, 0x12, 0x14, 0x0a,
	0x10, 0x44, 0x65, 0x61, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x45, 0x78, 0x63, 0x65, 0x65, 0x64, 0x65,
	0x64, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10,
	0x05, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x45, 0x78, 0x69, 0x73,
	0x74, 0x73, 0x10, 0x06, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x44, 0x65, 0x6e, 0x69, 0x65, 0x64, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x45, 0x78, 0x68, 0x61, 0x75, 0x73, 0x74, 0x65, 0x64, 0x10,
	0x08, 0x12, 0x16, 0x0a, 0x12, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x50, 0x72, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x09, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x62, 0x6f,
	0x72, 0x74, 0x65, 0x64, 0x10, 0x0a, 0x12, 0x0e, 0x0a, 0x0a, 0x4f, 0x75, 0x74, 0x4f, 0x66, 0x52,
	0x61, 0x6e, 0x67, 0x65, 0x10, 0x0b, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x6e, 0x69, 0x6d, 0x70, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x65, 0x64, 0x10, 0x0c, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x10, 0x0d, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x61, 0x76, 0x61,
	0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x10, 0x0e, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61,
	0x4c, 0x6f, 0x73, 0x73, 0x10, 0x0f, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x6e, 0x61, 0x75, 0x74, 0x68,
	0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x10, 0x10, 0x12, 0x13, 0x0a, 0x0f, 0x54,
	0x6f, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x10, 0x11,
	0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70,
	0x75, 0x62, 0x67, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x3b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errorpb_code_proto_rawDescOnce sync.Once
	file_errorpb_code_proto_rawDescData = file_errorpb_code_proto_rawDesc
)

func file_errorpb_code_proto_rawDescGZIP() []byte {
	file_errorpb_code_proto_rawDescOnce.Do(func() {
		file_errorpb_code_proto_rawDescData = protoimpl.X.CompressGZIP(file_errorpb_code_proto_rawDescData)
	})
	return file_errorpb_code_proto_rawDescData
}

var file_errorpb_code_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errorpb_code_proto_goTypes = []any{
	(Code)(0), // 0: errors.Code
}
var file_errorpb_code_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_errorpb_code_proto_init() }
func file_errorpb_code_proto_init() {
	if File_errorpb_code_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_errorpb_code_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errorpb_code_proto_goTypes,
		DependencyIndexes: file_errorpb_code_proto_depIdxs,
		EnumInfos:         file_errorpb_code_proto_enumTypes,
	}.Build()
	File_errorpb_code_proto = out.File
	file_errorpb_code_proto_rawDesc = nil
	file_errorpb_code_proto_goTypes = nil
	file_errorpb_code_proto_depIdxs = nil
}
