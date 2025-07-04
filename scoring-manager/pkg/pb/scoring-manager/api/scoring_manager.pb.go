// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v6.30.0--rc1
// source: api/scoring_manager.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ScoringStatus int32

const (
	ScoringStatus_INITIAL ScoringStatus = 0
	ScoringStatus_SUCCESS ScoringStatus = 1
	ScoringStatus_FAILED  ScoringStatus = 2
)

// Enum value maps for ScoringStatus.
var (
	ScoringStatus_name = map[int32]string{
		0: "INITIAL",
		1: "SUCCESS",
		2: "FAILED",
	}
	ScoringStatus_value = map[string]int32{
		"INITIAL": 0,
		"SUCCESS": 1,
		"FAILED":  2,
	}
)

func (x ScoringStatus) Enum() *ScoringStatus {
	p := new(ScoringStatus)
	*p = x
	return p
}

func (x ScoringStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ScoringStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_scoring_manager_proto_enumTypes[0].Descriptor()
}

func (ScoringStatus) Type() protoreflect.EnumType {
	return &file_api_scoring_manager_proto_enumTypes[0]
}

func (x ScoringStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ScoringStatus.Descriptor instead.
func (ScoringStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_scoring_manager_proto_rawDescGZIP(), []int{0}
}

type StartScoringRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartScoringRequest) Reset() {
	*x = StartScoringRequest{}
	mi := &file_api_scoring_manager_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartScoringRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartScoringRequest) ProtoMessage() {}

func (x *StartScoringRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_scoring_manager_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartScoringRequest.ProtoReflect.Descriptor instead.
func (*StartScoringRequest) Descriptor() ([]byte, []int) {
	return file_api_scoring_manager_proto_rawDescGZIP(), []int{0}
}

func (x *StartScoringRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type StartScoringResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ApplicationId int64                  `protobuf:"varint,1,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartScoringResponse) Reset() {
	*x = StartScoringResponse{}
	mi := &file_api_scoring_manager_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartScoringResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartScoringResponse) ProtoMessage() {}

func (x *StartScoringResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_scoring_manager_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartScoringResponse.ProtoReflect.Descriptor instead.
func (*StartScoringResponse) Descriptor() ([]byte, []int) {
	return file_api_scoring_manager_proto_rawDescGZIP(), []int{1}
}

func (x *StartScoringResponse) GetApplicationId() int64 {
	if x != nil {
		return x.ApplicationId
	}
	return 0
}

type GetStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ApplicationId int64                  `protobuf:"varint,1,opt,name=application_id,json=applicationId,proto3" json:"application_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStatusRequest) Reset() {
	*x = GetStatusRequest{}
	mi := &file_api_scoring_manager_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatusRequest) ProtoMessage() {}

func (x *GetStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_scoring_manager_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatusRequest.ProtoReflect.Descriptor instead.
func (*GetStatusRequest) Descriptor() ([]byte, []int) {
	return file_api_scoring_manager_proto_rawDescGZIP(), []int{2}
}

func (x *GetStatusRequest) GetApplicationId() int64 {
	if x != nil {
		return x.ApplicationId
	}
	return 0
}

type GetStatusResponse struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Status  ScoringStatus          `protobuf:"varint,1,opt,name=status,proto3,enum=scoring_manager.ScoringStatus" json:"status,omitempty"`
	Scoring int32                  `protobuf:"varint,2,opt,name=scoring,proto3" json:"scoring,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStatusResponse) Reset() {
	*x = GetStatusResponse{}
	mi := &file_api_scoring_manager_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatusResponse) ProtoMessage() {}

func (x *GetStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_scoring_manager_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatusResponse.ProtoReflect.Descriptor instead.
func (*GetStatusResponse) Descriptor() ([]byte, []int) {
	return file_api_scoring_manager_proto_rawDescGZIP(), []int{3}
}

func (x *GetStatusResponse) GetStatus() ScoringStatus {
	if x != nil {
		return x.Status
	}
	return ScoringStatus_INITIAL
}

func (x *GetStatusResponse) GetScoring() int32 {
	if x != nil {
		return x.Scoring
	}
	return 0
}

var File_api_scoring_manager_proto protoreflect.FileDescriptor

var file_api_scoring_manager_proto_rawDesc = string([]byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x73, 0x63, 0x6f,
	0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22, 0x31, 0x0a, 0x13,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x3d, 0x0a, 0x14, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x39,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x65, 0x0a, 0x11, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e,
	0x2e, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x53, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x2a, 0x35, 0x0a, 0x0d, 0x53, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x49, 0x54, 0x49, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x0b,
	0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x32, 0xc8, 0x01, 0x0a, 0x15, 0x53, 0x63, 0x6f, 0x72,
	0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x5b, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x69, 0x6e,
	0x67, 0x12, 0x24, 0x2e, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e,
	0x67, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53,
	0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x52,
	0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x2e, 0x73, 0x63,
	0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22,
	0x2e, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x5d, 0x5a, 0x5b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x4b, 0x6c, 0x61, 0x73, 0x73, 0x6e, 0x61, 0x79, 0x61, 0x41, 0x66, 0x72, 0x6f, 0x64, 0x69,
	0x74, 0x61, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2d, 0x75, 0x73, 0x65, 0x72, 0x2d, 0x73,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2d, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62,
	0x3b, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_scoring_manager_proto_rawDescOnce sync.Once
	file_api_scoring_manager_proto_rawDescData []byte
)

func file_api_scoring_manager_proto_rawDescGZIP() []byte {
	file_api_scoring_manager_proto_rawDescOnce.Do(func() {
		file_api_scoring_manager_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_scoring_manager_proto_rawDesc), len(file_api_scoring_manager_proto_rawDesc)))
	})
	return file_api_scoring_manager_proto_rawDescData
}

var file_api_scoring_manager_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_scoring_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_scoring_manager_proto_goTypes = []any{
	(ScoringStatus)(0),           // 0: scoring_manager.ScoringStatus
	(*StartScoringRequest)(nil),  // 1: scoring_manager.StartScoringRequest
	(*StartScoringResponse)(nil), // 2: scoring_manager.StartScoringResponse
	(*GetStatusRequest)(nil),     // 3: scoring_manager.GetStatusRequest
	(*GetStatusResponse)(nil),    // 4: scoring_manager.GetStatusResponse
}
var file_api_scoring_manager_proto_depIdxs = []int32{
	0, // 0: scoring_manager.GetStatusResponse.status:type_name -> scoring_manager.ScoringStatus
	1, // 1: scoring_manager.ScoringManagerService.StartScoring:input_type -> scoring_manager.StartScoringRequest
	3, // 2: scoring_manager.ScoringManagerService.GetStatus:input_type -> scoring_manager.GetStatusRequest
	2, // 3: scoring_manager.ScoringManagerService.StartScoring:output_type -> scoring_manager.StartScoringResponse
	4, // 4: scoring_manager.ScoringManagerService.GetStatus:output_type -> scoring_manager.GetStatusResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_scoring_manager_proto_init() }
func file_api_scoring_manager_proto_init() {
	if File_api_scoring_manager_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_scoring_manager_proto_rawDesc), len(file_api_scoring_manager_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_scoring_manager_proto_goTypes,
		DependencyIndexes: file_api_scoring_manager_proto_depIdxs,
		EnumInfos:         file_api_scoring_manager_proto_enumTypes,
		MessageInfos:      file_api_scoring_manager_proto_msgTypes,
	}.Build()
	File_api_scoring_manager_proto = out.File
	file_api_scoring_manager_proto_goTypes = nil
	file_api_scoring_manager_proto_depIdxs = nil
}
