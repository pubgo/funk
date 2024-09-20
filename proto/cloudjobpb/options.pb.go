// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: cloudjobpb/options.proto

package cloudjobpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_cloudjobpb_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         10010,
		Name:          "lava.cloudjobs.job_name",
		Tag:           "bytes,10010,opt,name=job_name",
		Filename:      "cloudjobpb/options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         10011,
		Name:          "lava.cloudjobs.subject_name",
		Tag:           "bytes,10011,opt,name=subject_name",
		Filename:      "cloudjobpb/options.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// job name is same with config jobs consumers job name
	//
	// optional string job_name = 10010;
	E_JobName = &file_cloudjobpb_options_proto_extTypes[0]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// subject name is same with config jobs consumers
	// subject name sametime is topic name
	//
	// optional string subject_name = 10011;
	E_SubjectName = &file_cloudjobpb_options_proto_extTypes[1]
)

var File_cloudjobpb_options_proto protoreflect.FileDescriptor

var file_cloudjobpb_options_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x70, 0x62, 0x2f, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6c, 0x61, 0x76, 0x61,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3b, 0x0a, 0x08,
	0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x9a, 0x4e, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x3a, 0x42, 0x0a, 0x0c, 0x73, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x9b, 0x4e, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x33, 0x5a,
	0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67,
	0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x6a, 0x6f, 0x62, 0x70, 0x62, 0x3b, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_cloudjobpb_options_proto_goTypes = []any{
	(*descriptorpb.ServiceOptions)(nil), // 0: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 1: google.protobuf.MethodOptions
}
var file_cloudjobpb_options_proto_depIdxs = []int32{
	0, // 0: lava.cloudjobs.job_name:extendee -> google.protobuf.ServiceOptions
	1, // 1: lava.cloudjobs.subject_name:extendee -> google.protobuf.MethodOptions
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	0, // [0:2] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cloudjobpb_options_proto_init() }
func file_cloudjobpb_options_proto_init() {
	if File_cloudjobpb_options_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cloudjobpb_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_cloudjobpb_options_proto_goTypes,
		DependencyIndexes: file_cloudjobpb_options_proto_depIdxs,
		ExtensionInfos:    file_cloudjobpb_options_proto_extTypes,
	}.Build()
	File_cloudjobpb_options_proto = out.File
	file_cloudjobpb_options_proto_rawDesc = nil
	file_cloudjobpb_options_proto_goTypes = nil
	file_cloudjobpb_options_proto_depIdxs = nil
}
