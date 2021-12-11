// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: compute/msg.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

type MsgStoreCode struct {
	Sender github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=sender,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"sender,omitempty"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode []byte `protobuf:"bytes,2,opt,name=wasm_byte_code,json=wasmByteCode,proto3" json:"wasm_byte_code,omitempty"`
	// Source is a valid absolute HTTPS URI to the contract's source code, optional
	Source string `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	// Builder is a valid docker image name with tag, optional
	Builder string `protobuf:"bytes,4,opt,name=builder,proto3" json:"builder,omitempty"`
	// InstantiatePermission to apply on contract creation, optional
	//  AccessConfig InstantiatePermission = 5;
	ContractPeriod int64  `protobuf:"varint,5,opt,name=contract_period,json=contractPeriod,proto3" json:"contract_period,omitempty"`
	Title          string `protobuf:"bytes,6,opt,name=title,proto3" json:"title,omitempty"`
	Description    string `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	Instances      uint64 `protobuf:"varint,8,opt,name=instances,proto3" json:"instances,omitempty"`
}

func (m *MsgStoreCode) Reset()         { *m = MsgStoreCode{} }
func (m *MsgStoreCode) String() string { return proto.CompactTextString(m) }
func (*MsgStoreCode) ProtoMessage()    {}
func (*MsgStoreCode) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c8bfcb6dfda1556, []int{0}
}
func (m *MsgStoreCode) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgStoreCode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgStoreCode.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgStoreCode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgStoreCode.Merge(m, src)
}
func (m *MsgStoreCode) XXX_Size() int {
	return m.Size()
}
func (m *MsgStoreCode) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgStoreCode.DiscardUnknown(m)
}

var xxx_messageInfo_MsgStoreCode proto.InternalMessageInfo

type MsgInstantiateContract struct {
	Sender github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=sender,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"sender,omitempty"`
	// Admin is an optional address that can execute migrations
	//  bytes admin = 2 [(gogoproto.casttype) = "github.com/cosmos/cosmos-sdk/types.AccAddress"];
	CallbackCodeHash string                                   `protobuf:"bytes,2,opt,name=callback_code_hash,json=callbackCodeHash,proto3" json:"callback_code_hash,omitempty"`
	CodeID           uint64                                   `protobuf:"varint,3,opt,name=code_id,json=codeId,proto3" json:"code_id,omitempty"`
	ContractId       string                                   `protobuf:"bytes,4,opt,name=contract_id,json=contractId,proto3" json:"contract_id,omitempty"`
	InitMsg          []byte                                   `protobuf:"bytes,5,opt,name=init_msg,json=initMsg,proto3" json:"init_msg,omitempty"`
	LastMsg          []byte                                   `protobuf:"bytes,6,opt,name=last_msg,json=lastMsg,proto3" json:"last_msg,omitempty"`
	InitFunds        github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,7,rep,name=init_funds,json=initFunds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"init_funds"`
	CallbackSig      []byte                                   `protobuf:"bytes,8,opt,name=callback_sig,json=callbackSig,proto3" json:"callback_sig,omitempty"`
}

func (m *MsgInstantiateContract) Reset()         { *m = MsgInstantiateContract{} }
func (m *MsgInstantiateContract) String() string { return proto.CompactTextString(m) }
func (*MsgInstantiateContract) ProtoMessage()    {}
func (*MsgInstantiateContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c8bfcb6dfda1556, []int{1}
}
func (m *MsgInstantiateContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgInstantiateContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgInstantiateContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgInstantiateContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgInstantiateContract.Merge(m, src)
}
func (m *MsgInstantiateContract) XXX_Size() int {
	return m.Size()
}
func (m *MsgInstantiateContract) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgInstantiateContract.DiscardUnknown(m)
}

var xxx_messageInfo_MsgInstantiateContract proto.InternalMessageInfo

type MsgExecuteContract struct {
	Sender           github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=sender,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"sender,omitempty"`
	Contract         github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,2,opt,name=contract,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"contract,omitempty"`
	Msg              []byte                                        `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	CallbackCodeHash string                                        `protobuf:"bytes,4,opt,name=callback_code_hash,json=callbackCodeHash,proto3" json:"callback_code_hash,omitempty"`
	SentFunds        github_com_cosmos_cosmos_sdk_types.Coins      `protobuf:"bytes,5,rep,name=sent_funds,json=sentFunds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"sent_funds"`
	CallbackSig      []byte                                        `protobuf:"bytes,6,opt,name=callback_sig,json=callbackSig,proto3" json:"callback_sig,omitempty"`
}

func (m *MsgExecuteContract) Reset()         { *m = MsgExecuteContract{} }
func (m *MsgExecuteContract) String() string { return proto.CompactTextString(m) }
func (*MsgExecuteContract) ProtoMessage()    {}
func (*MsgExecuteContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c8bfcb6dfda1556, []int{2}
}
func (m *MsgExecuteContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgExecuteContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgExecuteContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgExecuteContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgExecuteContract.Merge(m, src)
}
func (m *MsgExecuteContract) XXX_Size() int {
	return m.Size()
}
func (m *MsgExecuteContract) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgExecuteContract.DiscardUnknown(m)
}

var xxx_messageInfo_MsgExecuteContract proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgStoreCode)(nil), "trst.x.compute.v1beta1.MsgStoreCode")
	proto.RegisterType((*MsgInstantiateContract)(nil), "trst.x.compute.v1beta1.MsgInstantiateContract")
	proto.RegisterType((*MsgExecuteContract)(nil), "trst.x.compute.v1beta1.MsgExecuteContract")
}

func init() { proto.RegisterFile("compute/msg.proto", fileDescriptor_5c8bfcb6dfda1556) }

var fileDescriptor_5c8bfcb6dfda1556 = []byte{
	// 640 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xb1, 0x6e, 0xdb, 0x3a,
	0x14, 0xb5, 0x22, 0x47, 0x4e, 0x68, 0x21, 0xc9, 0x23, 0x82, 0x40, 0x09, 0x1e, 0x24, 0x23, 0x6f,
	0x78, 0x1e, 0x1a, 0xa9, 0x4e, 0x81, 0x0e, 0xdd, 0xe2, 0xb4, 0x45, 0x8d, 0xc2, 0x40, 0xa1, 0x0c,
	0x05, 0xba, 0x18, 0x14, 0xc9, 0xca, 0x4c, 0x64, 0xd1, 0xd0, 0xa5, 0xdb, 0xe4, 0x0f, 0x3a, 0x76,
	0xe8, 0x07, 0x74, 0xee, 0xda, 0x9f, 0xc8, 0x98, 0xb1, 0x93, 0x5b, 0x38, 0x7f, 0xd1, 0xa9, 0x20,
	0x25, 0x25, 0x19, 0x1a, 0x20, 0x28, 0x9a, 0xc9, 0xbc, 0xe7, 0x5c, 0x1d, 0x92, 0xe7, 0x5c, 0x13,
	0xfd, 0x43, 0xe5, 0x64, 0x3a, 0x53, 0x3c, 0x9a, 0x40, 0x1a, 0x4e, 0x0b, 0xa9, 0x24, 0xde, 0x52,
	0x05, 0xa8, 0xf0, 0x34, 0xac, 0x98, 0xf0, 0x5d, 0x2f, 0xe1, 0x8a, 0xf4, 0x76, 0x36, 0x53, 0x99,
	0x4a, 0xd3, 0x12, 0xe9, 0x55, 0xd9, 0xbd, 0xe3, 0x53, 0x09, 0x13, 0x09, 0x51, 0x42, 0x80, 0x47,
	0x55, 0x6b, 0x44, 0xa5, 0xc8, 0x4b, 0x7e, 0xf7, 0x7c, 0x09, 0xb9, 0x43, 0x48, 0x8f, 0x94, 0x2c,
	0xf8, 0xa1, 0x64, 0x1c, 0x0f, 0x90, 0x03, 0x3c, 0x67, 0xbc, 0xf0, 0xac, 0x8e, 0xd5, 0x75, 0xfb,
	0xbd, 0x9f, 0xf3, 0x60, 0x2f, 0x15, 0x6a, 0x3c, 0x4b, 0xf4, 0x96, 0x51, 0xa5, 0x57, 0xfe, 0xec,
	0x01, 0x3b, 0x89, 0xd4, 0xd9, 0x94, 0x43, 0x78, 0x40, 0xe9, 0x01, 0x63, 0x05, 0x07, 0x88, 0x2b,
	0x01, 0xfc, 0x18, 0xad, 0xbd, 0x27, 0x30, 0x19, 0x25, 0x67, 0x8a, 0x8f, 0xa8, 0x64, 0xdc, 0x5b,
	0x32, 0x92, 0x1b, 0x8b, 0x79, 0xe0, 0xbe, 0x3e, 0x38, 0x1a, 0xf6, 0xcf, 0x94, 0xd9, 0x34, 0x76,
	0x75, 0x5f, 0x5d, 0xe1, 0x2d, 0xe4, 0x80, 0x9c, 0x15, 0x94, 0x7b, 0x76, 0xc7, 0xea, 0xae, 0xc6,
	0x55, 0x85, 0x3d, 0xd4, 0x4a, 0x66, 0x22, 0xd3, 0x67, 0x6b, 0x1a, 0xa2, 0x2e, 0xf1, 0xff, 0x68,
	0x9d, 0xca, 0x5c, 0x15, 0x84, 0xaa, 0xd1, 0x94, 0x17, 0x42, 0x32, 0x6f, 0xb9, 0x63, 0x75, 0xed,
	0x78, 0xad, 0x86, 0x5f, 0x19, 0x14, 0x6f, 0xa2, 0x65, 0x25, 0x54, 0xc6, 0x3d, 0xc7, 0x08, 0x94,
	0x05, 0xee, 0xa0, 0x36, 0xe3, 0x40, 0x0b, 0x31, 0x55, 0x42, 0xe6, 0x5e, 0xcb, 0x70, 0x37, 0x21,
	0xfc, 0x2f, 0x5a, 0x15, 0x39, 0x28, 0x92, 0x53, 0x0e, 0xde, 0x4a, 0xc7, 0xea, 0x36, 0xe3, 0x6b,
	0xe0, 0x49, 0xf3, 0xc3, 0xe7, 0xa0, 0xb1, 0xfb, 0xd5, 0x46, 0x5b, 0x43, 0x48, 0x07, 0x06, 0x56,
	0x82, 0xe8, 0xdb, 0x94, 0x9b, 0xff, 0x4d, 0x53, 0x1f, 0x20, 0x4c, 0x49, 0x96, 0x25, 0x84, 0x9e,
	0x18, 0x4f, 0x47, 0x63, 0x02, 0x63, 0x63, 0xec, 0x6a, 0xbc, 0x51, 0x33, 0xda, 0xc6, 0x17, 0x04,
	0xc6, 0xf8, 0x3f, 0xd4, 0x32, 0x4d, 0x82, 0x19, 0x2f, 0x9b, 0x7d, 0xb4, 0x98, 0x07, 0x8e, 0xa6,
	0x07, 0x4f, 0x63, 0x47, 0x53, 0x03, 0x86, 0x03, 0xd4, 0xbe, 0x72, 0x4f, 0xb0, 0xca, 0x5b, 0x54,
	0x43, 0x03, 0x86, 0xb7, 0xd1, 0x8a, 0xc8, 0x85, 0x1a, 0x4d, 0x20, 0x35, 0xbe, 0xba, 0x71, 0x4b,
	0xd7, 0x43, 0x48, 0x35, 0x95, 0x11, 0x28, 0x29, 0xa7, 0xa4, 0x74, 0xad, 0xa9, 0x63, 0x84, 0xcc,
	0x57, 0x6f, 0x67, 0x39, 0x03, 0xaf, 0xd5, 0xb1, 0xbb, 0xed, 0xfd, 0xed, 0xb0, 0xbc, 0x63, 0xa8,
	0xe7, 0xb1, 0x1e, 0xdd, 0xf0, 0x50, 0x8a, 0xbc, 0xff, 0xf0, 0x7c, 0x1e, 0x34, 0xbe, 0x7c, 0x0f,
	0xba, 0x77, 0xf0, 0x45, 0x7f, 0x00, 0x3a, 0x01, 0xa1, 0x9e, 0x6b, 0x75, 0xbc, 0x8f, 0xdc, 0x2b,
	0x57, 0x40, 0xa4, 0x26, 0x22, 0xb7, 0xbf, 0xbe, 0x98, 0x07, 0xed, 0xc3, 0x0a, 0x3f, 0x12, 0x69,
	0xdc, 0xa6, 0xd7, 0x45, 0x95, 0xda, 0x27, 0x1b, 0xe1, 0x21, 0xa4, 0xcf, 0x4e, 0x39, 0x9d, 0xdd,
	0x4f, 0x62, 0x43, 0xb4, 0x52, 0x7b, 0x59, 0xfd, 0x01, 0xfe, 0x40, 0xec, 0x4a, 0x02, 0x6f, 0x20,
	0x5b, 0x9b, 0x6d, 0x1b, 0xb3, 0xf5, 0xf2, 0x96, 0x91, 0x68, 0xde, 0x32, 0x12, 0xc7, 0x08, 0x01,
	0xcf, 0xeb, 0x58, 0x96, 0xef, 0x21, 0x16, 0x2d, 0xff, 0xfb, 0x58, 0x9c, 0xbb, 0xc6, 0xd2, 0x7f,
	0x79, 0xbe, 0xf0, 0xad, 0x8b, 0x85, 0x6f, 0xfd, 0x58, 0xf8, 0xd6, 0xc7, 0x4b, 0xbf, 0x71, 0x71,
	0xe9, 0x37, 0xbe, 0x5d, 0xfa, 0x8d, 0x37, 0xbd, 0x1b, 0x07, 0xd1, 0x4f, 0x61, 0x46, 0x12, 0x30,
	0x8b, 0xe8, 0x34, 0xaa, 0x5f, 0x4b, 0x91, 0x2b, 0x5e, 0xe4, 0x24, 0x2b, 0xcf, 0x95, 0x38, 0xe6,
	0xad, 0x7b, 0xf4, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xf5, 0xd2, 0x4e, 0xde, 0x4e, 0x05, 0x00, 0x00,
}

func (m *MsgStoreCode) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgStoreCode) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgStoreCode) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Instances != 0 {
		i = encodeVarintMsg(dAtA, i, uint64(m.Instances))
		i--
		dAtA[i] = 0x40
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0x32
	}
	if m.ContractPeriod != 0 {
		i = encodeVarintMsg(dAtA, i, uint64(m.ContractPeriod))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Builder) > 0 {
		i -= len(m.Builder)
		copy(dAtA[i:], m.Builder)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Builder)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Source) > 0 {
		i -= len(m.Source)
		copy(dAtA[i:], m.Source)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Source)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.WASMByteCode) > 0 {
		i -= len(m.WASMByteCode)
		copy(dAtA[i:], m.WASMByteCode)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.WASMByteCode)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgInstantiateContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgInstantiateContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgInstantiateContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.CallbackSig) > 0 {
		i -= len(m.CallbackSig)
		copy(dAtA[i:], m.CallbackSig)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.CallbackSig)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.InitFunds) > 0 {
		for iNdEx := len(m.InitFunds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.InitFunds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMsg(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.LastMsg) > 0 {
		i -= len(m.LastMsg)
		copy(dAtA[i:], m.LastMsg)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.LastMsg)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.InitMsg) > 0 {
		i -= len(m.InitMsg)
		copy(dAtA[i:], m.InitMsg)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.InitMsg)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.ContractId) > 0 {
		i -= len(m.ContractId)
		copy(dAtA[i:], m.ContractId)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.ContractId)))
		i--
		dAtA[i] = 0x22
	}
	if m.CodeID != 0 {
		i = encodeVarintMsg(dAtA, i, uint64(m.CodeID))
		i--
		dAtA[i] = 0x18
	}
	if len(m.CallbackCodeHash) > 0 {
		i -= len(m.CallbackCodeHash)
		copy(dAtA[i:], m.CallbackCodeHash)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.CallbackCodeHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgExecuteContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgExecuteContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgExecuteContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.CallbackSig) > 0 {
		i -= len(m.CallbackSig)
		copy(dAtA[i:], m.CallbackSig)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.CallbackSig)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.SentFunds) > 0 {
		for iNdEx := len(m.SentFunds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SentFunds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMsg(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.CallbackCodeHash) > 0 {
		i -= len(m.CallbackCodeHash)
		copy(dAtA[i:], m.CallbackCodeHash)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.CallbackCodeHash)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintMsg(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsg(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsg(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgStoreCode) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.WASMByteCode)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.Source)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.Builder)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if m.ContractPeriod != 0 {
		n += 1 + sovMsg(uint64(m.ContractPeriod))
	}
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if m.Instances != 0 {
		n += 1 + sovMsg(uint64(m.Instances))
	}
	return n
}

func (m *MsgInstantiateContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.CallbackCodeHash)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if m.CodeID != 0 {
		n += 1 + sovMsg(uint64(m.CodeID))
	}
	l = len(m.ContractId)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.InitMsg)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.LastMsg)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if len(m.InitFunds) > 0 {
		for _, e := range m.InitFunds {
			l = e.Size()
			n += 1 + l + sovMsg(uint64(l))
		}
	}
	l = len(m.CallbackSig)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	return n
}

func (m *MsgExecuteContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	l = len(m.CallbackCodeHash)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	if len(m.SentFunds) > 0 {
		for _, e := range m.SentFunds {
			l = e.Size()
			n += 1 + l + sovMsg(uint64(l))
		}
	}
	l = len(m.CallbackSig)
	if l > 0 {
		n += 1 + l + sovMsg(uint64(l))
	}
	return n
}

func sovMsg(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsg(x uint64) (n int) {
	return sovMsg(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgStoreCode) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgStoreCode: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgStoreCode: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = append(m.Sender[:0], dAtA[iNdEx:postIndex]...)
			if m.Sender == nil {
				m.Sender = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WASMByteCode", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WASMByteCode = append(m.WASMByteCode[:0], dAtA[iNdEx:postIndex]...)
			if m.WASMByteCode == nil {
				m.WASMByteCode = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Source = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Builder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Builder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractPeriod", wireType)
			}
			m.ContractPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ContractPeriod |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instances", wireType)
			}
			m.Instances = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Instances |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgInstantiateContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgInstantiateContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgInstantiateContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = append(m.Sender[:0], dAtA[iNdEx:postIndex]...)
			if m.Sender == nil {
				m.Sender = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallbackCodeHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CallbackCodeHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeID", wireType)
			}
			m.CodeID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CodeID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitMsg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InitMsg = append(m.InitMsg[:0], dAtA[iNdEx:postIndex]...)
			if m.InitMsg == nil {
				m.InitMsg = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastMsg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastMsg = append(m.LastMsg[:0], dAtA[iNdEx:postIndex]...)
			if m.LastMsg == nil {
				m.LastMsg = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitFunds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InitFunds = append(m.InitFunds, types.Coin{})
			if err := m.InitFunds[len(m.InitFunds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallbackSig", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CallbackSig = append(m.CallbackSig[:0], dAtA[iNdEx:postIndex]...)
			if m.CallbackSig == nil {
				m.CallbackSig = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgExecuteContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsg
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgExecuteContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgExecuteContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = append(m.Sender[:0], dAtA[iNdEx:postIndex]...)
			if m.Sender == nil {
				m.Sender = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = append(m.Contract[:0], dAtA[iNdEx:postIndex]...)
			if m.Contract == nil {
				m.Contract = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallbackCodeHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CallbackCodeHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SentFunds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SentFunds = append(m.SentFunds, types.Coin{})
			if err := m.SentFunds[len(m.SentFunds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallbackSig", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsg
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CallbackSig = append(m.CallbackSig[:0], dAtA[iNdEx:postIndex]...)
			if m.CallbackSig == nil {
				m.CallbackSig = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsg
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMsg(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsg
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsg
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMsg
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsg
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsg
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsg        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsg          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsg = fmt.Errorf("proto: unexpected end of group")
)
