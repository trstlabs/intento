// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: intento/claim/v1beta1/claim.proto

package types

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Action int32

const (
	ACTION_ACTION_LOCAL    Action = 0
	ACTION_ACTION_ICA      Action = 1
	ACTION_GOVERNANCE_VOTE Action = 2
	ACTION_DELEGATE_STAKE  Action = 3
)

var Action_name = map[int32]string{
	0: "ACTION_ACTION_LOCAL",
	1: "ACTION_ACTION_ICA",
	2: "ACTION_GOVERNANCE_VOTE",
	3: "ACTION_DELEGATE_STAKE",
}

var Action_value = map[string]int32{
	"ACTION_ACTION_LOCAL":    0,
	"ACTION_ACTION_ICA":      1,
	"ACTION_GOVERNANCE_VOTE": 2,
	"ACTION_DELEGATE_STAKE":  3,
}

func (x Action) String() string {
	return proto.EnumName(Action_name, int32(x))
}

func (Action) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c568070ade923d6a, []int{0}
}

// A Claim Records is the metadata of claim data per address
type ClaimRecord struct {
	// address of recipient
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	// maximum claimable amount for the address
	MaximumClaimableAmount types.Coin `protobuf:"bytes,2,opt,name=maximum_claimable_amount,json=maximumClaimableAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"maximum_claimable_amount" yaml:"maximum_claimable_amount"`
	// index of status array refers to action enum #
	Status []Status `protobuf:"bytes,3,rep,name=status,proto3" json:"status" yaml:"status"`
}

func (m *ClaimRecord) Reset()         { *m = ClaimRecord{} }
func (m *ClaimRecord) String() string { return proto.CompactTextString(m) }
func (*ClaimRecord) ProtoMessage()    {}
func (*ClaimRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_c568070ade923d6a, []int{0}
}
func (m *ClaimRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClaimRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClaimRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClaimRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClaimRecord.Merge(m, src)
}
func (m *ClaimRecord) XXX_Size() int {
	return m.Size()
}
func (m *ClaimRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_ClaimRecord.DiscardUnknown(m)
}

var xxx_messageInfo_ClaimRecord proto.InternalMessageInfo

func (m *ClaimRecord) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *ClaimRecord) GetMaximumClaimableAmount() types.Coin {
	if m != nil {
		return m.MaximumClaimableAmount
	}
	return types.Coin{}
}

func (m *ClaimRecord) GetStatus() []Status {
	if m != nil {
		return m.Status
	}
	return nil
}

// Status contains for an action if it is completed and claimed
type Status struct {
	// true if action is completed
	ActionCompleted bool `protobuf:"varint,1,opt,name=action_completed,json=actionCompleted,proto3" json:"action_completed,omitempty" yaml:"action_completed"`
	// true if action is completed
	// index refers to the 4 vesting periods for the given action
	VestingPeriodsCompleted []bool `protobuf:"varint,2,rep,packed,name=vesting_periods_completed,json=vestingPeriodsCompleted,proto3" json:"vesting_periods_completed,omitempty" yaml:"vesting_periods_completed"`
	// true if action is completed
	// index refers to the 4 vesting periods for the given action
	VestingPeriodsClaimed []bool `protobuf:"varint,3,rep,packed,name=vesting_periods_claimed,json=vestingPeriodsClaimed,proto3" json:"vesting_periods_claimed,omitempty" yaml:"vesting_periods_claimed"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_c568070ade923d6a, []int{1}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Status.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return m.Size()
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetActionCompleted() bool {
	if m != nil {
		return m.ActionCompleted
	}
	return false
}

func (m *Status) GetVestingPeriodsCompleted() []bool {
	if m != nil {
		return m.VestingPeriodsCompleted
	}
	return nil
}

func (m *Status) GetVestingPeriodsClaimed() []bool {
	if m != nil {
		return m.VestingPeriodsClaimed
	}
	return nil
}

type MsgClaimClaimable struct {
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (m *MsgClaimClaimable) Reset()         { *m = MsgClaimClaimable{} }
func (m *MsgClaimClaimable) String() string { return proto.CompactTextString(m) }
func (*MsgClaimClaimable) ProtoMessage()    {}
func (*MsgClaimClaimable) Descriptor() ([]byte, []int) {
	return fileDescriptor_c568070ade923d6a, []int{2}
}
func (m *MsgClaimClaimable) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimClaimable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimClaimable.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimClaimable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimClaimable.Merge(m, src)
}
func (m *MsgClaimClaimable) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimClaimable) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimClaimable.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimClaimable proto.InternalMessageInfo

func (m *MsgClaimClaimable) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

type MsgClaimClaimableResponse struct {
	// returned claimable amount for the address
	ClaimedAmount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=claimed_amount,json=claimedAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"claimed_amount" yaml:"claimed_amount"`
}

func (m *MsgClaimClaimableResponse) Reset()         { *m = MsgClaimClaimableResponse{} }
func (m *MsgClaimClaimableResponse) String() string { return proto.CompactTextString(m) }
func (*MsgClaimClaimableResponse) ProtoMessage()    {}
func (*MsgClaimClaimableResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c568070ade923d6a, []int{3}
}
func (m *MsgClaimClaimableResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimClaimableResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimClaimableResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimClaimableResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimClaimableResponse.Merge(m, src)
}
func (m *MsgClaimClaimableResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimClaimableResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimClaimableResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimClaimableResponse proto.InternalMessageInfo

func (m *MsgClaimClaimableResponse) GetClaimedAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.ClaimedAmount
	}
	return nil
}

func init() {
	proto.RegisterEnum("intento.claim.v1beta1.Action", Action_name, Action_value)
	proto.RegisterType((*ClaimRecord)(nil), "intento.claim.v1beta1.ClaimRecord")
	proto.RegisterType((*Status)(nil), "intento.claim.v1beta1.Status")
	proto.RegisterType((*MsgClaimClaimable)(nil), "intento.claim.v1beta1.MsgClaimClaimable")
	proto.RegisterType((*MsgClaimClaimableResponse)(nil), "intento.claim.v1beta1.MsgClaimClaimableResponse")
}

func init() { proto.RegisterFile("intento/claim/v1beta1/claim.proto", fileDescriptor_c568070ade923d6a) }

var fileDescriptor_c568070ade923d6a = []byte{
	// 676 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xbf, 0x6f, 0xd3, 0x40,
	0x14, 0x8e, 0x13, 0x08, 0xe5, 0xaa, 0x96, 0xf4, 0x20, 0xcd, 0x0f, 0x84, 0x1d, 0x2c, 0x24, 0x42,
	0x04, 0x36, 0x2d, 0x5b, 0x37, 0xdb, 0x98, 0xaa, 0x22, 0x6d, 0x90, 0x1b, 0x75, 0xe8, 0x62, 0x9c,
	0xf8, 0x64, 0x2c, 0x62, 0x5f, 0x94, 0xbb, 0x54, 0x2d, 0x13, 0x62, 0x42, 0x20, 0x21, 0xfe, 0x07,
	0x16, 0x60, 0xea, 0xc0, 0x1f, 0xd1, 0xb1, 0x23, 0x93, 0x41, 0xed, 0xd0, 0x3d, 0x7f, 0x01, 0xca,
	0xdd, 0x99, 0xfe, 0x08, 0x41, 0xb0, 0xe4, 0xf2, 0xde, 0xf7, 0xdd, 0x77, 0xef, 0x7d, 0xf7, 0x7c,
	0xe0, 0x76, 0x18, 0x53, 0x14, 0x53, 0xac, 0x77, 0x7b, 0x5e, 0x18, 0xe9, 0x3b, 0x4b, 0x1d, 0x44,
	0xbd, 0x25, 0x1e, 0x69, 0xfd, 0x01, 0xa6, 0x18, 0x16, 0x05, 0x45, 0xe3, 0x49, 0x41, 0xa9, 0xde,
	0x08, 0x70, 0x80, 0x19, 0x43, 0x1f, 0xff, 0xe3, 0xe4, 0xaa, 0xdc, 0xc5, 0x24, 0xc2, 0x44, 0xef,
	0x78, 0x04, 0x9d, 0xaa, 0xe1, 0x30, 0x16, 0x78, 0x49, 0xe0, 0x11, 0x09, 0xf4, 0x9d, 0xa5, 0xf1,
	0x22, 0x80, 0x05, 0x2f, 0x0a, 0x63, 0xac, 0xb3, 0x5f, 0x9e, 0x52, 0xbf, 0x65, 0xc1, 0xac, 0x35,
	0x3e, 0xd3, 0x41, 0x5d, 0x3c, 0xf0, 0xe1, 0x7d, 0x70, 0xc5, 0xf3, 0xfd, 0x01, 0x22, 0xa4, 0x2c,
	0xd5, 0xa4, 0xfa, 0x55, 0x13, 0x8e, 0x12, 0x65, 0x7e, 0xcf, 0x8b, 0x7a, 0x2b, 0xaa, 0x00, 0x54,
	0x27, 0xa5, 0xc0, 0xcf, 0x12, 0x28, 0x47, 0xde, 0x6e, 0x18, 0x0d, 0x23, 0x97, 0x55, 0xee, 0x75,
	0x7a, 0xc8, 0xf5, 0x22, 0x3c, 0x8c, 0x69, 0x39, 0x5b, 0x93, 0xea, 0xb3, 0xcb, 0x15, 0x8d, 0x57,
	0xa3, 0x8d, 0xab, 0x4d, 0x1b, 0xd3, 0x2c, 0x1c, 0xc6, 0xe6, 0xe6, 0x41, 0xa2, 0x64, 0x46, 0x89,
	0xa2, 0x70, 0xf9, 0x69, 0x42, 0xea, 0xd7, 0x1f, 0x4a, 0x3d, 0x08, 0xe9, 0x8b, 0x61, 0x47, 0xeb,
	0xe2, 0x48, 0x17, 0xdd, 0xf1, 0xe5, 0x01, 0xf1, 0x5f, 0xea, 0x74, 0xaf, 0x8f, 0x08, 0xd3, 0x24,
	0xce, 0xa2, 0x90, 0xb1, 0x52, 0x15, 0x83, 0x89, 0xc0, 0x26, 0xc8, 0x13, 0xea, 0xd1, 0x21, 0x29,
	0xe7, 0x6a, 0xb9, 0xfa, 0xec, 0xf2, 0x2d, 0xed, 0x8f, 0x96, 0x6b, 0x9b, 0x8c, 0x64, 0x16, 0x45,
	0x6d, 0x73, 0xbc, 0x36, 0xbe, 0x55, 0x75, 0x84, 0x86, 0xfa, 0x21, 0x0b, 0xf2, 0x9c, 0x09, 0x9f,
	0x80, 0x82, 0xd7, 0xa5, 0x21, 0x8e, 0xdd, 0x2e, 0x8e, 0xfa, 0x3d, 0x44, 0x91, 0xcf, 0xac, 0x9b,
	0x31, 0x6f, 0x8e, 0x12, 0xa5, 0x24, 0xac, 0xbb, 0xc0, 0x50, 0x9d, 0x6b, 0x3c, 0x65, 0xa5, 0x19,
	0xf8, 0x1c, 0x54, 0x76, 0x10, 0xa1, 0x61, 0x1c, 0xb8, 0x7d, 0x34, 0x08, 0xb1, 0x4f, 0xce, 0x08,
	0x66, 0x6b, 0xb9, 0xfa, 0x8c, 0x79, 0x67, 0x94, 0x28, 0x35, 0x2e, 0x38, 0x95, 0xaa, 0x3a, 0x25,
	0x81, 0x3d, 0xe3, 0xd0, 0xe9, 0x09, 0xdb, 0xa0, 0x34, 0xb1, 0x6d, 0xdc, 0x3b, 0xf2, 0x99, 0x27,
	0x33, 0xa6, 0x3a, 0x4a, 0x14, 0x79, 0x8a, 0x3e, 0x27, 0xaa, 0x4e, 0xf1, 0x82, 0xba, 0xc8, 0xb7,
	0xc1, 0xc2, 0x3a, 0x09, 0x58, 0xf4, 0xdb, 0x79, 0xb8, 0x08, 0xf2, 0x04, 0xc5, 0x3e, 0x1a, 0xf0,
	0x59, 0x72, 0x44, 0xb4, 0x72, 0xf7, 0xcd, 0xc9, 0x7e, 0x43, 0x04, 0xef, 0x4e, 0xf6, 0x1b, 0x25,
	0xfe, 0x61, 0x4c, 0x08, 0xa8, 0x5f, 0x24, 0x50, 0x99, 0xc8, 0x3a, 0x88, 0xf4, 0x71, 0x4c, 0x10,
	0x7c, 0x2f, 0x81, 0x79, 0x51, 0x57, 0x3a, 0x73, 0x12, 0xbb, 0xdb, 0xbf, 0xcc, 0xdc, 0x9a, 0xb8,
	0xd7, 0x22, 0x6f, 0xf3, 0xfc, 0xf6, 0xff, 0x9b, 0xb4, 0x39, 0xb1, 0x99, 0x0f, 0x58, 0x63, 0x08,
	0xf2, 0x06, 0xbb, 0x52, 0x58, 0x02, 0xd7, 0x0d, 0xab, 0xbd, 0xd6, 0xda, 0x70, 0xc5, 0xd2, 0x6c,
	0x59, 0x46, 0xb3, 0x90, 0x81, 0x45, 0xb0, 0x70, 0x1e, 0x58, 0xb3, 0x8c, 0x82, 0x04, 0xab, 0x60,
	0x51, 0xc4, 0xab, 0xad, 0x2d, 0xdb, 0xd9, 0x30, 0x36, 0x2c, 0xdb, 0xdd, 0x6a, 0xb5, 0xed, 0x42,
	0x16, 0x56, 0x40, 0x51, 0x60, 0x8f, 0xed, 0xa6, 0xbd, 0x6a, 0xb4, 0x6d, 0x77, 0xb3, 0x6d, 0x3c,
	0xb5, 0x0b, 0xb9, 0xea, 0xa5, 0xb7, 0x9f, 0xe4, 0xcc, 0xf2, 0x2b, 0x90, 0x5b, 0x27, 0x01, 0xec,
	0x81, 0xf9, 0x0b, 0xe6, 0xd7, 0xa7, 0x0c, 0xf8, 0x84, 0x9f, 0xd5, 0x87, 0xff, 0xca, 0x4c, 0x9d,
	0xaf, 0x5e, 0x7e, 0x7d, 0xb2, 0xdf, 0x90, 0x4c, 0xeb, 0xe0, 0x48, 0x96, 0x0e, 0x8f, 0x64, 0xe9,
	0xe7, 0x91, 0x2c, 0x7d, 0x3c, 0x96, 0x33, 0x87, 0xc7, 0x72, 0xe6, 0xfb, 0xb1, 0x9c, 0xd9, 0xbe,
	0x77, 0xc6, 0x45, 0x3a, 0x20, 0xb4, 0xe7, 0x75, 0x88, 0x9e, 0x3e, 0x83, 0xbb, 0xe2, 0x21, 0x64,
	0x66, 0x76, 0xf2, 0xec, 0x21, 0x7a, 0xf4, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x16, 0xa1, 0x57, 0xc3,
	0x26, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	ClaimClaimable(ctx context.Context, in *MsgClaimClaimable, opts ...grpc.CallOption) (*MsgClaimClaimableResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) ClaimClaimable(ctx context.Context, in *MsgClaimClaimable, opts ...grpc.CallOption) (*MsgClaimClaimableResponse, error) {
	out := new(MsgClaimClaimableResponse)
	err := c.cc.Invoke(ctx, "/intento.claim.v1beta1.Msg/ClaimClaimable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	ClaimClaimable(context.Context, *MsgClaimClaimable) (*MsgClaimClaimableResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) ClaimClaimable(ctx context.Context, req *MsgClaimClaimable) (*MsgClaimClaimableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimClaimable not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_ClaimClaimable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimClaimable)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimClaimable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/intento.claim.v1beta1.Msg/ClaimClaimable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimClaimable(ctx, req.(*MsgClaimClaimable))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_serviceDesc = _Msg_serviceDesc
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "intento.claim.v1beta1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClaimClaimable",
			Handler:    _Msg_ClaimClaimable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "intento/claim/v1beta1/claim.proto",
}

func (m *ClaimRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClaimRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClaimRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Status) > 0 {
		for iNdEx := len(m.Status) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Status[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintClaim(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size, err := m.MaximumClaimableAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintClaim(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintClaim(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Status) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.VestingPeriodsClaimed) > 0 {
		for iNdEx := len(m.VestingPeriodsClaimed) - 1; iNdEx >= 0; iNdEx-- {
			i--
			if m.VestingPeriodsClaimed[iNdEx] {
				dAtA[i] = 1
			} else {
				dAtA[i] = 0
			}
		}
		i = encodeVarintClaim(dAtA, i, uint64(len(m.VestingPeriodsClaimed)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.VestingPeriodsCompleted) > 0 {
		for iNdEx := len(m.VestingPeriodsCompleted) - 1; iNdEx >= 0; iNdEx-- {
			i--
			if m.VestingPeriodsCompleted[iNdEx] {
				dAtA[i] = 1
			} else {
				dAtA[i] = 0
			}
		}
		i = encodeVarintClaim(dAtA, i, uint64(len(m.VestingPeriodsCompleted)))
		i--
		dAtA[i] = 0x12
	}
	if m.ActionCompleted {
		i--
		if m.ActionCompleted {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MsgClaimClaimable) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimClaimable) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimClaimable) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintClaim(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgClaimClaimableResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimClaimableResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimClaimableResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ClaimedAmount) > 0 {
		for iNdEx := len(m.ClaimedAmount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ClaimedAmount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintClaim(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintClaim(dAtA []byte, offset int, v uint64) int {
	offset -= sovClaim(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClaimRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovClaim(uint64(l))
	}
	l = m.MaximumClaimableAmount.Size()
	n += 1 + l + sovClaim(uint64(l))
	if len(m.Status) > 0 {
		for _, e := range m.Status {
			l = e.Size()
			n += 1 + l + sovClaim(uint64(l))
		}
	}
	return n
}

func (m *Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActionCompleted {
		n += 2
	}
	if len(m.VestingPeriodsCompleted) > 0 {
		n += 1 + sovClaim(uint64(len(m.VestingPeriodsCompleted))) + len(m.VestingPeriodsCompleted)*1
	}
	if len(m.VestingPeriodsClaimed) > 0 {
		n += 1 + sovClaim(uint64(len(m.VestingPeriodsClaimed))) + len(m.VestingPeriodsClaimed)*1
	}
	return n
}

func (m *MsgClaimClaimable) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovClaim(uint64(l))
	}
	return n
}

func (m *MsgClaimClaimableResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ClaimedAmount) > 0 {
		for _, e := range m.ClaimedAmount {
			l = e.Size()
			n += 1 + l + sovClaim(uint64(l))
		}
	}
	return n
}

func sovClaim(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozClaim(x uint64) (n int) {
	return sovClaim(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClaimRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClaim
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
			return fmt.Errorf("proto: ClaimRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClaimRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClaim
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
				return ErrInvalidLengthClaim
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthClaim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaximumClaimableAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClaim
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
				return ErrInvalidLengthClaim
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthClaim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaximumClaimableAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClaim
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
				return ErrInvalidLengthClaim
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthClaim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = append(m.Status, Status{})
			if err := m.Status[len(m.Status)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClaim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthClaim
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
func (m *Status) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClaim
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
			return fmt.Errorf("proto: Status: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Status: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActionCompleted", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClaim
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.ActionCompleted = bool(v != 0)
		case 2:
			if wireType == 0 {
				var v int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClaim
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.VestingPeriodsCompleted = append(m.VestingPeriodsCompleted, bool(v != 0))
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClaim
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthClaim
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthClaim
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen
				if elementCount != 0 && len(m.VestingPeriodsCompleted) == 0 {
					m.VestingPeriodsCompleted = make([]bool, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowClaim
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.VestingPeriodsCompleted = append(m.VestingPeriodsCompleted, bool(v != 0))
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingPeriodsCompleted", wireType)
			}
		case 3:
			if wireType == 0 {
				var v int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClaim
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.VestingPeriodsClaimed = append(m.VestingPeriodsClaimed, bool(v != 0))
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClaim
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthClaim
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthClaim
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen
				if elementCount != 0 && len(m.VestingPeriodsClaimed) == 0 {
					m.VestingPeriodsClaimed = make([]bool, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowClaim
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.VestingPeriodsClaimed = append(m.VestingPeriodsClaimed, bool(v != 0))
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingPeriodsClaimed", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipClaim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthClaim
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
func (m *MsgClaimClaimable) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClaim
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
			return fmt.Errorf("proto: MsgClaimClaimable: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimClaimable: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClaim
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
				return ErrInvalidLengthClaim
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthClaim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClaim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthClaim
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
func (m *MsgClaimClaimableResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClaim
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
			return fmt.Errorf("proto: MsgClaimClaimableResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimClaimableResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimedAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClaim
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
				return ErrInvalidLengthClaim
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthClaim
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClaimedAmount = append(m.ClaimedAmount, types.Coin{})
			if err := m.ClaimedAmount[len(m.ClaimedAmount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClaim(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthClaim
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
func skipClaim(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowClaim
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
					return 0, ErrIntOverflowClaim
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
					return 0, ErrIntOverflowClaim
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
				return 0, ErrInvalidLengthClaim
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupClaim
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthClaim
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthClaim        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowClaim          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupClaim = fmt.Errorf("proto: unexpected end of group")
)
