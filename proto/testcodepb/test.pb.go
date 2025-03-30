// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.3
// source: testcodepb/test.proto

package testcodepb

import (
	_ "github.com/pubgo/funk/pkg/gen/cloudjobpb"
	_ "github.com/pubgo/funk/proto/errorpb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/descriptorpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type GenType int32

const (
	GenType_default   GenType = 0
	GenType_uuid      GenType = 1
	GenType_snowflake GenType = 2
	GenType_bigflake  GenType = 3
	GenType_shortid   GenType = 4
)

// Enum value maps for GenType.
var (
	GenType_name = map[int32]string{
		0: "default",
		1: "uuid",
		2: "snowflake",
		3: "bigflake",
		4: "shortid",
	}
	GenType_value = map[string]int32{
		"default":   0,
		"uuid":      1,
		"snowflake": 2,
		"bigflake":  3,
		"shortid":   4,
	}
)

func (x GenType) Enum() *GenType {
	p := new(GenType)
	*p = x
	return p
}

func (x GenType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GenType) Descriptor() protoreflect.EnumDescriptor {
	return file_testcodepb_test_proto_enumTypes[1].Descriptor()
}

func (GenType) Type() protoreflect.EnumType {
	return &file_testcodepb_test_proto_enumTypes[1]
}

func (x GenType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GenType.Descriptor instead.
func (GenType) EnumDescriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{1}
}

type User_Role_Type int32

const (
	// 租户
	User_Role_tenant User_Role_Type = 0
	// 安保
	User_Role_guard User_Role_Type = 1
	// 管理者
	User_Role_manager User_Role_Type = 2
	// 后台管理员
	User_Role_admin User_Role_Type = 3
)

// Enum value maps for User_Role_Type.
var (
	User_Role_Type_name = map[int32]string{
		0: "tenant",
		1: "guard",
		2: "manager",
		3: "admin",
	}
	User_Role_Type_value = map[string]int32{
		"tenant":  0,
		"guard":   1,
		"manager": 2,
		"admin":   3,
	}
)

func (x User_Role_Type) Enum() *User_Role_Type {
	p := new(User_Role_Type)
	*p = x
	return p
}

func (x User_Role_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (User_Role_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_testcodepb_test_proto_enumTypes[2].Descriptor()
}

func (User_Role_Type) Type() protoreflect.EnumType {
	return &file_testcodepb_test_proto_enumTypes[2]
}

func (x User_Role_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use User_Role_Type.Descriptor instead.
func (User_Role_Type) EnumDescriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{0, 0, 0}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_testcodepb_test_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{0}
}

type GenerateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the unique id generated
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// the type of id generated
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *GenerateResponse) Reset() {
	*x = GenerateResponse{}
	mi := &file_testcodepb_test_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateResponse) ProtoMessage() {}

func (x *GenerateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateResponse.ProtoReflect.Descriptor instead.
func (*GenerateResponse) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{1}
}

func (x *GenerateResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GenerateResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type DoProxyEventReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DoProxyEventReq) Reset() {
	*x = DoProxyEventReq{}
	mi := &file_testcodepb_test_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DoProxyEventReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoProxyEventReq) ProtoMessage() {}

func (x *DoProxyEventReq) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoProxyEventReq.ProtoReflect.Descriptor instead.
func (*DoProxyEventReq) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{2}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_testcodepb_test_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{3}
}

type UploadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string             `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	File     *httpbody.HttpBody `protobuf:"bytes,2,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	mi := &file_testcodepb_test_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileRequest.ProtoReflect.Descriptor instead.
func (*UploadFileRequest) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{4}
}

func (x *UploadFileRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *UploadFileRequest) GetFile() *httpbody.HttpBody {
	if x != nil {
		return x.File
	}
	return nil
}

type ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SessionId string   `protobuf:"bytes,3,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Msg       *Message `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	mi := &file_testcodepb_test_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{5}
}

func (x *ChatMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChatMessage) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *ChatMessage) GetMsg() *Message {
	if x != nil {
		return x.Msg
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Text string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_testcodepb_test_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{6}
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Message) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

// Generate a unique ID. Defaults to uuid.
type GenerateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// type of id e.g uuid, shortid, snowflake (64 bit), bigflake (128 bit)
	Type GenType `protobuf:"varint,1,opt,name=type,proto3,enum=demo.test.v1.GenType" json:"type,omitempty"`
}

func (x *GenerateRequest) Reset() {
	*x = GenerateRequest{}
	mi := &file_testcodepb_test_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GenerateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateRequest) ProtoMessage() {}

func (x *GenerateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateRequest.ProtoReflect.Descriptor instead.
func (*GenerateRequest) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{7}
}

func (x *GenerateRequest) GetType() GenType {
	if x != nil {
		return x.Type
	}
	return GenType_default
}

// List the types of IDs available. No query params needed.
type TypesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	NameId uint64 `protobuf:"varint,2,opt,name=name_id,json=nameId,proto3" json:"name_id,omitempty"`
	Hello  string `protobuf:"bytes,3,opt,name=hello,proto3" json:"hello,omitempty"`
}

func (x *TypesRequest) Reset() {
	*x = TypesRequest{}
	mi := &file_testcodepb_test_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TypesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypesRequest) ProtoMessage() {}

func (x *TypesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypesRequest.ProtoReflect.Descriptor instead.
func (*TypesRequest) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{8}
}

func (x *TypesRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TypesRequest) GetNameId() uint64 {
	if x != nil {
		return x.NameId
	}
	return 0
}

func (x *TypesRequest) GetHello() string {
	if x != nil {
		return x.Hello
	}
	return ""
}

// TypesResponse 返回值类型
type TypesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Types []string `protobuf:"bytes,1,rep,name=types,proto3" json:"types,omitempty"`
}

func (x *TypesResponse) Reset() {
	*x = TypesResponse{}
	mi := &file_testcodepb_test_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TypesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypesResponse) ProtoMessage() {}

func (x *TypesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypesResponse.ProtoReflect.Descriptor instead.
func (*TypesResponse) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{9}
}

func (x *TypesResponse) GetTypes() []string {
	if x != nil {
		return x.Types
	}
	return nil
}

type User_Role struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *User_Role) Reset() {
	*x = User_Role{}
	mi := &file_testcodepb_test_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User_Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User_Role) ProtoMessage() {}

func (x *User_Role) ProtoReflect() protoreflect.Message {
	mi := &file_testcodepb_test_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User_Role.ProtoReflect.Descriptor instead.
func (*User_Role) Descriptor() ([]byte, []int) {
	return file_testcodepb_test_proto_rawDescGZIP(), []int{0, 0}
}

var File_testcodepb_test_proto protoreflect.FileDescriptor

var file_testcodepb_test_proto_rawDesc = []byte{
	0x0a, 0x15, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x2f, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65,
	0x73, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x70, 0x62,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x62, 0x6f, 0x64, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x6a, 0x6f, 0x62, 0x73, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x3d,
	0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x35, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a,
	0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x10, 0x03, 0x22, 0x36, 0x0a,
	0x10, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x44, 0x6f, 0x50, 0x72, 0x6f, 0x78, 0x79,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x59, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x48, 0x74,
	0x74, 0x70, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x69, 0x0a, 0x0b,
	0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x27,
	0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x64, 0x65,
	0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x31, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x3c, 0x0a, 0x0f, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x64, 0x65,
	0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x51, 0x0a, 0x0c, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e,
	0x61, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x22, 0x25, 0x0a, 0x0d, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2a, 0xb4, 0x01, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f,
	0x4b, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10,
	0xa0, 0x8d, 0x06, 0x1a, 0x06, 0x92, 0xea, 0x30, 0x02, 0x08, 0x05, 0x12, 0x15, 0x0a, 0x07, 0x55,
	0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0xa1, 0x8d, 0x06, 0x1a, 0x06, 0x92, 0xea, 0x30, 0x02,
	0x08, 0x05, 0x12, 0x14, 0x0a, 0x06, 0x44, 0x62, 0x43, 0x6f, 0x6e, 0x6e, 0x10, 0xa3, 0x8d, 0x06,
	0x1a, 0x06, 0x92, 0xea, 0x30, 0x02, 0x08, 0x0d, 0x12, 0x11, 0x0a, 0x0b, 0x55, 0x6e, 0x6b, 0x6e,
	0x6f, 0x77, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x10, 0xa4, 0x8d, 0x06, 0x12, 0x3c, 0x0a, 0x0a, 0x43,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x43, 0x6f, 0x64, 0x65, 0x10, 0xa5, 0x8d, 0x06, 0x1a, 0x2a, 0x92,
	0xea, 0x30, 0x26, 0x12, 0x12, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x20, 0x6d, 0x73, 0x67, 0x1a, 0x10, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x63, 0x75,
	0x73, 0x74, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x1a, 0x0e, 0x8a, 0xea, 0x30, 0x0a, 0x08,
	0x01, 0x10, 0x0d, 0x1a, 0x04, 0x54, 0x65, 0x73, 0x74, 0x2a, 0x4a, 0x0a, 0x07, 0x47, 0x65, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x10,
	0x00, 0x12, 0x08, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x73,
	0x6e, 0x6f, 0x77, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x62, 0x69,
	0x67, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x69, 0x64, 0x10, 0x04, 0x32, 0xcc, 0x07, 0x0a, 0x02, 0x49, 0x64, 0x12, 0x65, 0x0a, 0x08,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x3a,
	0x01, 0x2a, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x64, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x12, 0x5e, 0x0a, 0x0a, 0x54, 0x79, 0x70, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x1a, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e,
	0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x64, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x31, 0x30, 0x01, 0x12, 0x67, 0x0a, 0x05, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x64,
	0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x12, 0x1d, 0x2f,
	0x76, 0x31, 0x2f, 0x69, 0x64, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d,
	0x65, 0x7d, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x6d, 0x0a, 0x08,
	0x50, 0x75, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x3a, 0x01, 0x2a, 0x1a, 0x1d, 0x2f, 0x76,
	0x31, 0x2f, 0x69, 0x64, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65,
	0x7d, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x57, 0x0a, 0x04, 0x43,
	0x68, 0x61, 0x74, 0x12, 0x19, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x19,
	0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68,
	0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0f, 0x3a, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x08, 0x2f, 0x77, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x74,
	0x28, 0x01, 0x30, 0x01, 0x12, 0x43, 0x0a, 0x05, 0x43, 0x68, 0x61, 0x74, 0x31, 0x12, 0x19, 0x2e,
	0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x19, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e,
	0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x68, 0x0a, 0x0e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1f, 0x2e, 0x64, 0x65,
	0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x42, 0x6f,
	0x64, 0x79, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x04, 0x66, 0x69, 0x6c, 0x65,
	0x22, 0x11, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61,
	0x6d, 0x65, 0x7d, 0x12, 0x4b, 0x0a, 0x07, 0x44, 0x6f, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x12, 0x13,
	0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x13, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0d, 0x12, 0x0b, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x12, 0x5b, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x45, 0x78, 0x65, 0x63, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x6f, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x12, 0xda, 0xf1, 0x04, 0x0e, 0x67,
	0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x65, 0x78, 0x65, 0x63, 0x12, 0x5b, 0x0a,
	0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x1d, 0x2e,
	0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x6f, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x14, 0xda, 0xf1, 0x04, 0x10, 0x67, 0x69, 0x64, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x18, 0xca, 0x41, 0x0e, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x3a, 0x38, 0x30, 0x38, 0x30, 0xd2, 0xf1, 0x04,
	0x03, 0x67, 0x69, 0x64, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x2f, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x3b, 0x74,
	0x65, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
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

var file_testcodepb_test_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_testcodepb_test_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_testcodepb_test_proto_goTypes = []any{
	(Code)(0),                 // 0: demo.test.v1.Code
	(GenType)(0),              // 1: demo.test.v1.GenType
	(User_Role_Type)(0),       // 2: demo.test.v1.User.Role.Type
	(*User)(nil),              // 3: demo.test.v1.User
	(*GenerateResponse)(nil),  // 4: demo.test.v1.GenerateResponse
	(*DoProxyEventReq)(nil),   // 5: demo.test.v1.DoProxyEventReq
	(*Empty)(nil),             // 6: demo.test.v1.Empty
	(*UploadFileRequest)(nil), // 7: demo.test.v1.UploadFileRequest
	(*ChatMessage)(nil),       // 8: demo.test.v1.ChatMessage
	(*Message)(nil),           // 9: demo.test.v1.Message
	(*GenerateRequest)(nil),   // 10: demo.test.v1.GenerateRequest
	(*TypesRequest)(nil),      // 11: demo.test.v1.TypesRequest
	(*TypesResponse)(nil),     // 12: demo.test.v1.TypesResponse
	(*User_Role)(nil),         // 13: demo.test.v1.User.Role
	(*httpbody.HttpBody)(nil), // 14: google.api.HttpBody
	(*emptypb.Empty)(nil),     // 15: google.protobuf.Empty
}
var file_testcodepb_test_proto_depIdxs = []int32{
	14, // 0: demo.test.v1.UploadFileRequest.file:type_name -> google.api.HttpBody
	9,  // 1: demo.test.v1.ChatMessage.msg:type_name -> demo.test.v1.Message
	1,  // 2: demo.test.v1.GenerateRequest.type:type_name -> demo.test.v1.GenType
	10, // 3: demo.test.v1.Id.Generate:input_type -> demo.test.v1.GenerateRequest
	11, // 4: demo.test.v1.Id.TypeStream:input_type -> demo.test.v1.TypesRequest
	11, // 5: demo.test.v1.Id.Types:input_type -> demo.test.v1.TypesRequest
	11, // 6: demo.test.v1.Id.PutTypes:input_type -> demo.test.v1.TypesRequest
	8,  // 7: demo.test.v1.Id.Chat:input_type -> demo.test.v1.ChatMessage
	8,  // 8: demo.test.v1.Id.Chat1:input_type -> demo.test.v1.ChatMessage
	7,  // 9: demo.test.v1.Id.UploadDownload:input_type -> demo.test.v1.UploadFileRequest
	6,  // 10: demo.test.v1.Id.DoProxy:input_type -> demo.test.v1.Empty
	5,  // 11: demo.test.v1.Id.ProxyExecEvent:input_type -> demo.test.v1.DoProxyEventReq
	5,  // 12: demo.test.v1.Id.EventChanged:input_type -> demo.test.v1.DoProxyEventReq
	4,  // 13: demo.test.v1.Id.Generate:output_type -> demo.test.v1.GenerateResponse
	12, // 14: demo.test.v1.Id.TypeStream:output_type -> demo.test.v1.TypesResponse
	12, // 15: demo.test.v1.Id.Types:output_type -> demo.test.v1.TypesResponse
	12, // 16: demo.test.v1.Id.PutTypes:output_type -> demo.test.v1.TypesResponse
	8,  // 17: demo.test.v1.Id.Chat:output_type -> demo.test.v1.ChatMessage
	8,  // 18: demo.test.v1.Id.Chat1:output_type -> demo.test.v1.ChatMessage
	14, // 19: demo.test.v1.Id.UploadDownload:output_type -> google.api.HttpBody
	15, // 20: demo.test.v1.Id.DoProxy:output_type -> google.protobuf.Empty
	15, // 21: demo.test.v1.Id.ProxyExecEvent:output_type -> google.protobuf.Empty
	15, // 22: demo.test.v1.Id.EventChanged:output_type -> google.protobuf.Empty
	13, // [13:23] is the sub-list for method output_type
	3,  // [3:13] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
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
			NumEnums:      3,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_testcodepb_test_proto_goTypes,
		DependencyIndexes: file_testcodepb_test_proto_depIdxs,
		EnumInfos:         file_testcodepb_test_proto_enumTypes,
		MessageInfos:      file_testcodepb_test_proto_msgTypes,
	}.Build()
	File_testcodepb_test_proto = out.File
	file_testcodepb_test_proto_rawDesc = nil
	file_testcodepb_test_proto_goTypes = nil
	file_testcodepb_test_proto_depIdxs = nil
}
