// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: msgtype.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

import proto "google.golang.org/protobuf/proto"

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MSG_TYPE int32

const (
	MSG_TYPE__ERROR                MSG_TYPE = 0  //Error
	MSG_TYPE__Login                MSG_TYPE = 1  //SdkLogin
	MSG_TYPE__LoginResult          MSG_TYPE = 2  //LoginResult
	MSG_TYPE__PlayStart            MSG_TYPE = 3  //PlayStart
	MSG_TYPE__PlayStartResult      MSG_TYPE = 4  //PlayStartResult
	MSG_TYPE__PlayEnd              MSG_TYPE = 5  //PlayEnd
	MSG_TYPE__PlayEndResult        MSG_TYPE = 6  //PlayEndResult
	MSG_TYPE__Report               MSG_TYPE = 7  //Report
	MSG_TYPE__ReportResult         MSG_TYPE = 8  //ReportResult
	MSG_TYPE__GetRank              MSG_TYPE = 9  //GetRank
	MSG_TYPE__GetRankResult        MSG_TYPE = 10 //GetRankResult
	MSG_TYPE__NotifyNewAudience    MSG_TYPE = 11 //NotifyNewAudience
	MSG_TYPE__NotifyAudienceAction MSG_TYPE = 12 //NotifyAudienceAction
	MSG_TYPE__Ping                 MSG_TYPE = 21 //Ping
	MSG_TYPE__Pong                 MSG_TYPE = 22 //Pong
)

// Enum value maps for MSG_TYPE.
var (
	MSG_TYPE_name = map[int32]string{
		0:  "_ERROR",
		1:  "_Login",
		2:  "_LoginResult",
		3:  "_PlayStart",
		4:  "_PlayStartResult",
		5:  "_PlayEnd",
		6:  "_PlayEndResult",
		7:  "_Report",
		8:  "_ReportResult",
		9:  "_GetRank",
		10: "_GetRankResult",
		11: "_NotifyNewAudience",
		12: "_NotifyAudienceAction",
		21: "_Ping",
		22: "_Pong",
	}
	MSG_TYPE_value = map[string]int32{
		"_ERROR":                0,
		"_Login":                1,
		"_LoginResult":          2,
		"_PlayStart":            3,
		"_PlayStartResult":      4,
		"_PlayEnd":              5,
		"_PlayEndResult":        6,
		"_Report":               7,
		"_ReportResult":         8,
		"_GetRank":              9,
		"_GetRankResult":        10,
		"_NotifyNewAudience":    11,
		"_NotifyAudienceAction": 12,
		"_Ping":                 21,
		"_Pong":                 22,
	}
)

func (x MSG_TYPE) Enum() *MSG_TYPE {
	p := new(MSG_TYPE)
	*p = x
	return p
}

func (x MSG_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MSG_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_msgtype_proto_enumTypes[0].Descriptor()
}

func (MSG_TYPE) Type() protoreflect.EnumType {
	return &file_msgtype_proto_enumTypes[0]
}

func (x MSG_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MSG_TYPE.Descriptor instead.
func (MSG_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_msgtype_proto_rawDescGZIP(), []int{0}
}

type ERROR_CODE int32

const (
	ERROR_CODE_SUCCESS         ERROR_CODE = 0
	ERROR_CODE_FAIL            ERROR_CODE = 1
	ERROR_CODE_INVALID_APPID   ERROR_CODE = 2
	ERROR_CODE_INVALID_TOKEN   ERROR_CODE = 3
	ERROR_CODE_GAME_IS_RUNNING ERROR_CODE = 4
	ERROR_CODE_GAME_IS_STOPPED ERROR_CODE = 5
)

// Enum value maps for ERROR_CODE.
var (
	ERROR_CODE_name = map[int32]string{
		0: "SUCCESS",
		1: "FAIL",
		2: "INVALID_APPID",
		3: "INVALID_TOKEN",
		4: "GAME_IS_RUNNING",
		5: "GAME_IS_STOPPED",
	}
	ERROR_CODE_value = map[string]int32{
		"SUCCESS":         0,
		"FAIL":            1,
		"INVALID_APPID":   2,
		"INVALID_TOKEN":   3,
		"GAME_IS_RUNNING": 4,
		"GAME_IS_STOPPED": 5,
	}
)

func (x ERROR_CODE) Enum() *ERROR_CODE {
	p := new(ERROR_CODE)
	*p = x
	return p
}

func (x ERROR_CODE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ERROR_CODE) Descriptor() protoreflect.EnumDescriptor {
	return file_msgtype_proto_enumTypes[1].Descriptor()
}

func (ERROR_CODE) Type() protoreflect.EnumType {
	return &file_msgtype_proto_enumTypes[1]
}

func (x ERROR_CODE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ERROR_CODE.Descriptor instead.
func (ERROR_CODE) EnumDescriptor() ([]byte, []int) {
	return file_msgtype_proto_rawDescGZIP(), []int{1}
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgtype_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (x *Error) FromDB(data []byte) error {
	return proto.Unmarshal(data, x)
}

func (x *Error) ToDB() ([]byte, error) {
	return proto.Marshal(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_msgtype_proto_msgTypes[0]
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
	return file_msgtype_proto_rawDescGZIP(), []int{0}
}

var File_msgtype_proto protoreflect.FileDescriptor

var file_msgtype_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x73, 0x67, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x2a, 0x87, 0x02, 0x0a, 0x08, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x0a,
	0x0a, 0x06, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x5f, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x5f, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x5f, 0x50, 0x6c, 0x61,
	0x79, 0x53, 0x74, 0x61, 0x72, 0x74, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x5f, 0x50, 0x6c, 0x61,
	0x79, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x04, 0x12, 0x0c,
	0x0a, 0x08, 0x5f, 0x50, 0x6c, 0x61, 0x79, 0x45, 0x6e, 0x64, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e,
	0x5f, 0x50, 0x6c, 0x61, 0x79, 0x45, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x06,
	0x12, 0x0b, 0x0a, 0x07, 0x5f, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x10, 0x07, 0x12, 0x11, 0x0a,
	0x0d, 0x5f, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x08,
	0x12, 0x0c, 0x0a, 0x08, 0x5f, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x10, 0x09, 0x12, 0x12,
	0x0a, 0x0e, 0x5f, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x10, 0x0a, 0x12, 0x16, 0x0a, 0x12, 0x5f, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x4e, 0x65, 0x77,
	0x41, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x10, 0x0b, 0x12, 0x19, 0x0a, 0x15, 0x5f, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x79, 0x41, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x10, 0x0c, 0x12, 0x09, 0x0a, 0x05, 0x5f, 0x50, 0x69, 0x6e, 0x67, 0x10, 0x15,
	0x12, 0x09, 0x0a, 0x05, 0x5f, 0x50, 0x6f, 0x6e, 0x67, 0x10, 0x16, 0x2a, 0x73, 0x0a, 0x0a, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43,
	0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x01,
	0x12, 0x11, 0x0a, 0x0d, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41, 0x50, 0x50, 0x49,
	0x44, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x54,
	0x4f, 0x4b, 0x45, 0x4e, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x49,
	0x53, 0x5f, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x13, 0x0a, 0x0f, 0x47,
	0x41, 0x4d, 0x45, 0x5f, 0x49, 0x53, 0x5f, 0x53, 0x54, 0x4f, 0x50, 0x50, 0x45, 0x44, 0x10, 0x05,
	0x42, 0x05, 0x5a, 0x03, 0x70, 0x62, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msgtype_proto_rawDescOnce sync.Once
	file_msgtype_proto_rawDescData = file_msgtype_proto_rawDesc
)

func file_msgtype_proto_rawDescGZIP() []byte {
	file_msgtype_proto_rawDescOnce.Do(func() {
		file_msgtype_proto_rawDescData = protoimpl.X.CompressGZIP(file_msgtype_proto_rawDescData)
	})
	return file_msgtype_proto_rawDescData
}

var file_msgtype_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_msgtype_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_msgtype_proto_goTypes = []interface{}{
	(MSG_TYPE)(0),   // 0: message.MSG_TYPE
	(ERROR_CODE)(0), // 1: message.ERROR_CODE
	(*Error)(nil),   // 2: message.Error
}
var file_msgtype_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_msgtype_proto_init() }
func file_msgtype_proto_init() {
	if File_msgtype_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msgtype_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_msgtype_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msgtype_proto_goTypes,
		DependencyIndexes: file_msgtype_proto_depIdxs,
		EnumInfos:         file_msgtype_proto_enumTypes,
		MessageInfos:      file_msgtype_proto_msgTypes,
	}.Build()
	File_msgtype_proto = out.File
	file_msgtype_proto_rawDesc = nil
	file_msgtype_proto_goTypes = nil
	file_msgtype_proto_depIdxs = nil
}