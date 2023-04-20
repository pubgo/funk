// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: errorpb/errors.proto

package errorpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg    string            `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Detail string            `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
	Stack  string            `protobuf:"bytes,3,opt,name=stack,proto3" json:"stack,omitempty"`
	Tags   map[string]string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ErrMsg) Reset() {
	*x = ErrMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errorpb_errors_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrMsg) ProtoMessage() {}

func (x *ErrMsg) ProtoReflect() protoreflect.Message {
	mi := &file_errorpb_errors_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrMsg.ProtoReflect.Descriptor instead.
func (*ErrMsg) Descriptor() ([]byte, []int) {
	return file_errorpb_errors_proto_rawDescGZIP(), []int{0}
}

func (x *ErrMsg) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *ErrMsg) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *ErrMsg) GetStack() string {
	if x != nil {
		return x.Stack
	}
	return ""
}

func (x *ErrMsg) GetTags() map[string]string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type ErrCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   Code   `protobuf:"varint,2,opt,name=code,proto3,enum=status.Code" json:"code,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Reason string `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *ErrCode) Reset() {
	*x = ErrCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errorpb_errors_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrCode) ProtoMessage() {}

func (x *ErrCode) ProtoReflect() protoreflect.Message {
	mi := &file_errorpb_errors_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrCode.ProtoReflect.Descriptor instead.
func (*ErrCode) Descriptor() ([]byte, []int) {
	return file_errorpb_errors_proto_rawDescGZIP(), []int{1}
}

func (x *ErrCode) GetCode() Code {
	if x != nil {
		return x.Code
	}
	return Code_OK
}

func (x *ErrCode) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ErrCode) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type ErrTrace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Operation string `protobuf:"bytes,3,opt,name=operation,proto3" json:"operation,omitempty"`
	Service   string `protobuf:"bytes,4,opt,name=service,proto3" json:"service,omitempty"`
	Version   string `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *ErrTrace) Reset() {
	*x = ErrTrace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errorpb_errors_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrTrace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrTrace) ProtoMessage() {}

func (x *ErrTrace) ProtoReflect() protoreflect.Message {
	mi := &file_errorpb_errors_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrTrace.ProtoReflect.Descriptor instead.
func (*ErrTrace) Descriptor() ([]byte, []int) {
	return file_errorpb_errors_proto_rawDescGZIP(), []int{2}
}

func (x *ErrTrace) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ErrTrace) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *ErrTrace) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *ErrTrace) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code  *ErrCode  `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Trace *ErrTrace `protobuf:"bytes,2,opt,name=trace,proto3" json:"trace,omitempty"`
	Msg   *ErrMsg   `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errorpb_errors_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_errorpb_errors_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_errorpb_errors_proto_rawDescGZIP(), []int{3}
}

func (x *Error) GetCode() *ErrCode {
	if x != nil {
		return x.Code
	}
	return nil
}

func (x *Error) GetTrace() *ErrTrace {
	if x != nil {
		return x.Trace
	}
	return nil
}

func (x *Error) GetMsg() *ErrMsg {
	if x != nil {
		return x.Msg
	}
	return nil
}

var File_errorpb_errors_proto protoreflect.FileDescriptor

var file_errorpb_errors_proto_rawDesc = []byte{
	0x0a, 0x14, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x1a, 0x14,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xaf, 0x01, 0x0a, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x61,
	0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x2e, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x2e, 0x54, 0x61, 0x67, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x1a, 0x37, 0x0a, 0x09, 0x54, 0x61, 0x67, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x5b, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x6c,
	0x0a, 0x08, 0x45, 0x72, 0x72, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x76, 0x0a, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x23, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x45, 0x72, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x05, 0x74, 0x72, 0x61,
	0x63, 0x65, 0x12, 0x20, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x45, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62, 0x3b, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errorpb_errors_proto_rawDescOnce sync.Once
	file_errorpb_errors_proto_rawDescData = file_errorpb_errors_proto_rawDesc
)

func file_errorpb_errors_proto_rawDescGZIP() []byte {
	file_errorpb_errors_proto_rawDescOnce.Do(func() {
		file_errorpb_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_errorpb_errors_proto_rawDescData)
	})
	return file_errorpb_errors_proto_rawDescData
}

var file_errorpb_errors_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_errorpb_errors_proto_goTypes = []interface{}{
	(*ErrMsg)(nil),   // 0: errors.ErrMsg
	(*ErrCode)(nil),  // 1: errors.ErrCode
	(*ErrTrace)(nil), // 2: errors.ErrTrace
	(*Error)(nil),    // 3: errors.Error
	nil,              // 4: errors.ErrMsg.TagsEntry
	(Code)(0),        // 5: status.Code
}
var file_errorpb_errors_proto_depIdxs = []int32{
	4, // 0: errors.ErrMsg.tags:type_name -> errors.ErrMsg.TagsEntry
	5, // 1: errors.ErrCode.code:type_name -> status.Code
	1, // 2: errors.Error.code:type_name -> errors.ErrCode
	2, // 3: errors.Error.trace:type_name -> errors.ErrTrace
	0, // 4: errors.Error.msg:type_name -> errors.ErrMsg
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_errorpb_errors_proto_init() }
func file_errorpb_errors_proto_init() {
	if File_errorpb_errors_proto != nil {
		return
	}
	file_errorpb_status_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_errorpb_errors_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrMsg); i {
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
		file_errorpb_errors_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrCode); i {
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
		file_errorpb_errors_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrTrace); i {
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
		file_errorpb_errors_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
			RawDescriptor: file_errorpb_errors_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errorpb_errors_proto_goTypes,
		DependencyIndexes: file_errorpb_errors_proto_depIdxs,
		MessageInfos:      file_errorpb_errors_proto_msgTypes,
	}.Build()
	File_errorpb_errors_proto = out.File
	file_errorpb_errors_proto_rawDesc = nil
	file_errorpb_errors_proto_goTypes = nil
	file_errorpb_errors_proto_depIdxs = nil
}
