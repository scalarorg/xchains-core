// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: axelar/multisig/exported/v1beta1/types.proto

package exported

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MultisigState int32

const (
	NonExistent MultisigState = 0
	Pending     MultisigState = 1
	Completed   MultisigState = 2
)

var MultisigState_name = map[int32]string{
	0: "MULTISIG_STATE_UNSPECIFIED",
	1: "MULTISIG_STATE_PENDING",
	2: "MULTISIG_STATE_COMPLETED",
}

var MultisigState_value = map[string]int32{
	"MULTISIG_STATE_UNSPECIFIED": 0,
	"MULTISIG_STATE_PENDING":     1,
	"MULTISIG_STATE_COMPLETED":   2,
}

func (x MultisigState) String() string {
	return proto.EnumName(MultisigState_name, int32(x))
}

func (MultisigState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b14433678c926388, []int{0}
}

type KeyState int32

const (
	Inactive KeyState = 0
	Assigned KeyState = 1
	Active   KeyState = 2
)

var KeyState_name = map[int32]string{
	0: "KEY_STATE_UNSPECIFIED",
	1: "KEY_STATE_ASSIGNED",
	2: "KEY_STATE_ACTIVE",
}

var KeyState_value = map[string]int32{
	"KEY_STATE_UNSPECIFIED": 0,
	"KEY_STATE_ASSIGNED":    1,
	"KEY_STATE_ACTIVE":      2,
}

func (x KeyState) String() string {
	return proto.EnumName(KeyState_name, int32(x))
}

func (KeyState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b14433678c926388, []int{1}
}

func init() {
	proto.RegisterEnum("axelar.multisig.exported.v1beta1.MultisigState", MultisigState_name, MultisigState_value)
	proto.RegisterEnum("axelar.multisig.exported.v1beta1.KeyState", KeyState_name, KeyState_value)
}

func init() {
	proto.RegisterFile("axelar/multisig/exported/v1beta1/types.proto", fileDescriptor_b14433678c926388)
}

var fileDescriptor_b14433678c926388 = []byte{
	// 378 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0xd1, 0xc1, 0x8a, 0xda, 0x40,
	0x1c, 0x06, 0xf0, 0x8c, 0x07, 0x6b, 0xc7, 0x4a, 0x43, 0x68, 0x4b, 0xc9, 0x61, 0xc8, 0xa1, 0x20,
	0xd8, 0x36, 0x83, 0xf4, 0xd2, 0x6b, 0x9a, 0x4c, 0x65, 0x50, 0xd3, 0x40, 0xa2, 0xd0, 0x5e, 0x24,
	0x9a, 0x21, 0x0d, 0xd5, 0x4c, 0x48, 0x46, 0x1b, 0x1f, 0xa0, 0x50, 0x72, 0xea, 0x0b, 0x04, 0x0a,
	0xed, 0xa1, 0x8f, 0xe2, 0xd1, 0x63, 0x8f, 0xbb, 0xfa, 0x22, 0xcb, 0x1a, 0xc5, 0xc5, 0xdd, 0xdb,
	0x0c, 0xfc, 0xbe, 0x3f, 0x1f, 0x7c, 0xf0, 0x8d, 0x9f, 0xb3, 0xb9, 0x9f, 0xe2, 0xc5, 0x72, 0x2e,
	0xa2, 0x2c, 0x0a, 0x31, 0xcb, 0x13, 0x9e, 0x0a, 0x16, 0xe0, 0x55, 0x77, 0xca, 0x84, 0xdf, 0xc5,
	0x62, 0x9d, 0xb0, 0x4c, 0x4f, 0x52, 0x2e, 0xb8, 0xa2, 0x55, 0x5a, 0x3f, 0x69, 0xfd, 0xa4, 0xf5,
	0xa3, 0x56, 0x9f, 0x85, 0x3c, 0xe4, 0x07, 0x8c, 0x6f, 0x5f, 0x55, 0xae, 0xf3, 0x1b, 0xc0, 0xd6,
	0xf0, 0x98, 0x71, 0x85, 0x2f, 0x98, 0x82, 0xa1, 0x3a, 0x1c, 0x0d, 0x3c, 0xea, 0xd2, 0xde, 0xc4,
	0xf5, 0x0c, 0x8f, 0x4c, 0x46, 0xb6, 0xeb, 0x10, 0x93, 0x7e, 0xa4, 0xc4, 0x92, 0x25, 0xf5, 0x69,
	0x51, 0x6a, 0x4d, 0x9b, 0xc7, 0x24, 0x8f, 0x32, 0xc1, 0x62, 0xa1, 0xb4, 0xe1, 0x8b, 0x8b, 0x80,
	0x43, 0x6c, 0x8b, 0xda, 0x3d, 0x19, 0xa8, 0xcd, 0xa2, 0xd4, 0x1e, 0x39, 0x2c, 0x0e, 0xa2, 0x38,
	0x54, 0x5e, 0xc3, 0x97, 0x17, 0xd0, 0xfc, 0x34, 0x74, 0x06, 0xc4, 0x23, 0x96, 0x5c, 0x53, 0x5b,
	0x45, 0xa9, 0x3d, 0x36, 0xf9, 0x22, 0x99, 0x33, 0xc1, 0x02, 0xb5, 0xf1, 0xf3, 0x0f, 0x92, 0xfe,
	0xfd, 0x45, 0xa0, 0xf3, 0x03, 0xc0, 0x46, 0x9f, 0xad, 0xab, 0x76, 0x6d, 0xf8, 0xbc, 0x4f, 0x3e,
	0x3f, 0x58, 0xec, 0x49, 0x51, 0x6a, 0x0d, 0x1a, 0xfb, 0x33, 0x11, 0xad, 0x98, 0xf2, 0x0a, 0x2a,
	0x67, 0x68, 0xb8, 0x2e, 0xed, 0xd9, 0xc4, 0x92, 0x41, 0xa5, 0x8c, 0x2c, 0x8b, 0xc2, 0x98, 0x05,
	0x8a, 0x06, 0xe5, 0x3b, 0xca, 0xf4, 0xe8, 0x98, 0xc8, 0x35, 0x15, 0x16, 0xa5, 0x56, 0x37, 0x0e,
	0x77, 0xce, 0x3d, 0x3e, 0x8c, 0x37, 0xd7, 0x48, 0xda, 0xec, 0x10, 0xd8, 0xee, 0x10, 0xb8, 0xda,
	0x21, 0xf0, 0x6b, 0x8f, 0xa4, 0xed, 0x1e, 0x49, 0xff, 0xf7, 0x48, 0xfa, 0xf2, 0x3e, 0x8c, 0xc4,
	0xd7, 0xe5, 0x54, 0x9f, 0xf1, 0x05, 0xae, 0xb6, 0x88, 0x99, 0xf8, 0xce, 0xd3, 0x6f, 0xc7, 0xdf,
	0xdb, 0x19, 0x4f, 0x19, 0xce, 0xef, 0xcf, 0x39, 0xad, 0x1f, 0x96, 0x78, 0x77, 0x13, 0x00, 0x00,
	0xff, 0xff, 0xf2, 0x0f, 0x7f, 0x2c, 0xf1, 0x01, 0x00, 0x00,
}