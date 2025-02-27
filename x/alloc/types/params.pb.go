// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: intento/alloc/v1beta1/params.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// WeightedAddress defines an address with a weight.
type WeightedAddress struct {
	Address string                      `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	Weight  cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=weight,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight" yaml:"weight"`
}

func (m *WeightedAddress) Reset()         { *m = WeightedAddress{} }
func (m *WeightedAddress) String() string { return proto.CompactTextString(m) }
func (*WeightedAddress) ProtoMessage()    {}
func (*WeightedAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5691a07ae281fe, []int{0}
}
func (m *WeightedAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WeightedAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WeightedAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WeightedAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeightedAddress.Merge(m, src)
}
func (m *WeightedAddress) XXX_Size() int {
	return m.Size()
}
func (m *WeightedAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_WeightedAddress.DiscardUnknown(m)
}

var xxx_messageInfo_WeightedAddress proto.InternalMessageInfo

func (m *WeightedAddress) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// DistributionProportions defines the proportion that each bucket  receives.
type DistributionProportions struct {
	DeveloperRewards  cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=developer_rewards,json=developerRewards,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"developer_rewards" yaml:"developer_rewards"`
	CommunityPool     cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=community_pool,json=communityPool,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"community_pool" yaml:"community_pool"`
	RelayerIncentives cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=relayer_incentives,json=relayerIncentives,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"relayer_incentives" yaml:"community_pool"`
}

func (m *DistributionProportions) Reset()         { *m = DistributionProportions{} }
func (m *DistributionProportions) String() string { return proto.CompactTextString(m) }
func (*DistributionProportions) ProtoMessage()    {}
func (*DistributionProportions) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5691a07ae281fe, []int{1}
}
func (m *DistributionProportions) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DistributionProportions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DistributionProportions.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DistributionProportions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DistributionProportions.Merge(m, src)
}
func (m *DistributionProportions) XXX_Size() int {
	return m.Size()
}
func (m *DistributionProportions) XXX_DiscardUnknown() {
	xxx_messageInfo_DistributionProportions.DiscardUnknown(m)
}

var xxx_messageInfo_DistributionProportions proto.InternalMessageInfo

// Params defines the parameters for the alloc module.
type Params struct {
	// distribution_proportions defines the proportion of the minted denom
	DistributionProportions DistributionProportions `protobuf:"bytes,1,opt,name=distribution_proportions,json=distributionProportions,proto3" json:"distribution_proportions"`
	// addresses to receive developer rewards
	WeightedDeveloperRewardsReceivers []WeightedAddress `protobuf:"bytes,2,rep,name=weighted_developer_rewards_receivers,json=weightedDeveloperRewardsReceivers,proto3" json:"weighted_developer_rewards_receivers" yaml:"developer_rewards_receiver"`
	// SupplementAmount is the amount to be supplemented from the pool on top of
	// newly minted coins.
	SupplementAmount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,4,rep,name=supplement_amount,json=supplementAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"supplement_amount" yaml:"supplement_amount"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_bf5691a07ae281fe, []int{2}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetDistributionProportions() DistributionProportions {
	if m != nil {
		return m.DistributionProportions
	}
	return DistributionProportions{}
}

func (m *Params) GetWeightedDeveloperRewardsReceivers() []WeightedAddress {
	if m != nil {
		return m.WeightedDeveloperRewardsReceivers
	}
	return nil
}

func (m *Params) GetSupplementAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.SupplementAmount
	}
	return nil
}

func init() {
	proto.RegisterType((*WeightedAddress)(nil), "intento.alloc.v1beta1.WeightedAddress")
	proto.RegisterType((*DistributionProportions)(nil), "intento.alloc.v1beta1.DistributionProportions")
	proto.RegisterType((*Params)(nil), "intento.alloc.v1beta1.Params")
}

func init() {
	proto.RegisterFile("intento/alloc/v1beta1/params.proto", fileDescriptor_bf5691a07ae281fe)
}

var fileDescriptor_bf5691a07ae281fe = []byte{
	// 612 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xc1, 0x6b, 0xd4, 0x4e,
	0x18, 0xdd, 0x74, 0x7f, 0xf4, 0x47, 0xa7, 0xb6, 0xba, 0xc1, 0xd2, 0xb4, 0x42, 0xd2, 0x06, 0x91,
	0x5a, 0x6c, 0x42, 0xeb, 0xcd, 0x5b, 0xb7, 0x0b, 0x22, 0x08, 0x96, 0xbd, 0x08, 0x5e, 0xc2, 0x24,
	0x19, 0xb6, 0x43, 0x93, 0x7c, 0x61, 0x66, 0x76, 0xeb, 0x5e, 0xf4, 0x5f, 0x10, 0x8f, 0x9e, 0xbc,
	0x29, 0x3d, 0xf5, 0xe0, 0x1f, 0xd1, 0x63, 0x11, 0x04, 0xf1, 0xb0, 0xca, 0xee, 0xa1, 0xf7, 0xfd,
	0x0b, 0x24, 0x33, 0x93, 0x58, 0xbb, 0xbb, 0xd0, 0x8b, 0x97, 0xdd, 0x49, 0xbe, 0x97, 0xf7, 0xde,
	0xcc, 0xf7, 0xe6, 0x43, 0x2e, 0xcd, 0x04, 0xc9, 0x04, 0xf8, 0x38, 0x49, 0x20, 0xf2, 0x7b, 0xbb,
	0x21, 0x11, 0x78, 0xd7, 0xcf, 0x31, 0xc3, 0x29, 0xf7, 0x72, 0x06, 0x02, 0xcc, 0x15, 0x8d, 0xf1,
	0x24, 0xc6, 0xd3, 0x98, 0xf5, 0xbb, 0x1d, 0xe8, 0x80, 0x44, 0xf8, 0xc5, 0x4a, 0x81, 0xd7, 0xd7,
	0x22, 0xe0, 0x29, 0xf0, 0x40, 0x15, 0xd4, 0x83, 0x2e, 0x35, 0x70, 0x4a, 0x33, 0xf0, 0xe5, 0xaf,
	0x7e, 0x65, 0x2b, 0x80, 0x1f, 0x62, 0x4e, 0x2a, 0xf1, 0x08, 0x68, 0xa6, 0xea, 0xee, 0x47, 0x03,
	0xdd, 0x7e, 0x49, 0x68, 0xe7, 0x48, 0x90, 0x78, 0x3f, 0x8e, 0x19, 0xe1, 0xdc, 0x7c, 0x84, 0xfe,
	0xc7, 0x6a, 0x69, 0x19, 0x1b, 0xc6, 0xd6, 0x42, 0xd3, 0x1c, 0x0f, 0x9c, 0xe5, 0x3e, 0x4e, 0x93,
	0x27, 0xae, 0x2e, 0xb8, 0xed, 0x12, 0x62, 0x06, 0x68, 0xfe, 0x44, 0x12, 0x58, 0x73, 0x12, 0xfc,
	0xf4, 0x7c, 0xe0, 0xd4, 0x7e, 0x0c, 0x9c, 0x7b, 0x4a, 0x99, 0xc7, 0xc7, 0x1e, 0x05, 0x3f, 0xc5,
	0xe2, 0xc8, 0x7b, 0x4e, 0x3a, 0x38, 0xea, 0xb7, 0x48, 0x34, 0x1e, 0x38, 0x4b, 0x8a, 0x4f, 0x7d,
	0xea, 0x7e, 0xfd, 0xb2, 0x83, 0xf4, 0x56, 0x5a, 0x24, 0xfa, 0x7c, 0x79, 0xb6, 0x6d, 0xb4, 0x35,
	0xad, 0xfb, 0xbe, 0x8e, 0x56, 0x5b, 0x94, 0x0b, 0x46, 0xc3, 0xae, 0xa0, 0x90, 0x1d, 0x32, 0xc8,
	0x81, 0x15, 0x2b, 0x6e, 0xbe, 0x45, 0x8d, 0x98, 0xf4, 0x48, 0x02, 0x39, 0x61, 0x01, 0x23, 0x27,
	0x98, 0xc5, 0xa5, 0xe9, 0xf6, 0xcd, 0x7c, 0x58, 0xca, 0xc7, 0x04, 0xcb, 0x54, 0x4b, 0x77, 0x2a,
	0x58, 0x5b, 0xa1, 0xcc, 0x1e, 0x5a, 0x8e, 0x20, 0x4d, 0xbb, 0x19, 0x15, 0xfd, 0x20, 0x07, 0x48,
	0xf4, 0x29, 0xbc, 0xb8, 0x99, 0xfa, 0x8a, 0x52, 0xff, 0x9b, 0x62, 0xaa, 0xf4, 0x52, 0x85, 0x39,
	0x04, 0x48, 0xcc, 0x37, 0xc8, 0x64, 0x24, 0xc1, 0x7d, 0xc2, 0x02, 0x9a, 0x45, 0x24, 0x13, 0xb4,
	0x47, 0xb8, 0x55, 0xff, 0x37, 0xda, 0x0d, 0x2d, 0xf5, 0xac, 0x52, 0x72, 0xbf, 0xd5, 0xd1, 0xfc,
	0xa1, 0xcc, 0xb0, 0x29, 0x90, 0x15, 0x5f, 0x69, 0x4f, 0x11, 0xcc, 0xb2, 0x3f, 0xb2, 0x15, 0x8b,
	0x7b, 0x9e, 0x37, 0x35, 0xe0, 0xde, 0x8c, 0xae, 0x36, 0x17, 0x8a, 0x0d, 0x28, 0xe9, 0xd5, 0x78,
	0x46, 0xe7, 0x3f, 0x19, 0xe8, 0xfe, 0x89, 0x0e, 0x6e, 0x30, 0xd1, 0xbd, 0x80, 0x91, 0x88, 0xd0,
	0x1e, 0x61, 0xdc, 0x9a, 0xdb, 0xa8, 0x6f, 0x2d, 0xee, 0x3d, 0x98, 0x61, 0xe1, 0x5a, 0xf6, 0x9b,
	0x5e, 0x21, 0x3d, 0x1e, 0x38, 0x9b, 0x33, 0x62, 0x51, 0x11, 0xbb, 0xca, 0xdf, 0x66, 0xe9, 0xa1,
	0x75, 0x2d, 0x19, 0xed, 0xd2, 0x80, 0x79, 0x6a, 0xa0, 0x06, 0xef, 0xe6, 0x79, 0x42, 0x52, 0x92,
	0x89, 0x00, 0xa7, 0xd0, 0xcd, 0x84, 0xf5, 0x9f, 0xb4, 0xb5, 0xe6, 0xe9, 0xb3, 0x2e, 0xee, 0x67,
	0x65, 0xea, 0x00, 0x68, 0xd6, 0x0c, 0xb5, 0x13, 0x1d, 0xd0, 0x09, 0x06, 0xf7, 0xf4, 0xa7, 0xb3,
	0xd5, 0xa1, 0xe2, 0xa8, 0x1b, 0x7a, 0x11, 0xa4, 0x7a, 0x12, 0xe8, 0xbf, 0x1d, 0x1e, 0x1f, 0xfb,
	0xa2, 0x9f, 0x13, 0x2e, 0xc9, 0xf8, 0x87, 0xcb, 0xb3, 0xed, 0x5b, 0x89, 0x6c, 0x7e, 0x50, 0xdc,
	0x7e, 0xae, 0xf3, 0xfc, 0x87, 0x75, 0x5f, 0x92, 0x36, 0x0f, 0xce, 0x87, 0xb6, 0x71, 0x31, 0xb4,
	0x8d, 0x5f, 0x43, 0xdb, 0x78, 0x37, 0xb2, 0x6b, 0x17, 0x23, 0xbb, 0xf6, 0x7d, 0x64, 0xd7, 0x5e,
	0x3d, 0xbc, 0xa2, 0x25, 0x18, 0x17, 0x09, 0x0e, 0xb9, 0x5f, 0x0e, 0xb7, 0xd7, 0x7a, 0xbc, 0x49,
	0xc9, 0x70, 0x5e, 0xce, 0x96, 0xc7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xa0, 0xe8, 0xbe, 0x08,
	0xfc, 0x04, 0x00, 0x00,
}

func (m *WeightedAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WeightedAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WeightedAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Weight.Size()
		i -= size
		if _, err := m.Weight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DistributionProportions) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DistributionProportions) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DistributionProportions) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.RelayerIncentives.Size()
		i -= size
		if _, err := m.RelayerIncentives.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.CommunityPool.Size()
		i -= size
		if _, err := m.CommunityPool.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.DeveloperRewards.Size()
		i -= size
		if _, err := m.DeveloperRewards.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SupplementAmount) > 0 {
		for iNdEx := len(m.SupplementAmount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SupplementAmount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.WeightedDeveloperRewardsReceivers) > 0 {
		for iNdEx := len(m.WeightedDeveloperRewardsReceivers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.WeightedDeveloperRewardsReceivers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.DistributionProportions.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WeightedAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.Weight.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func (m *DistributionProportions) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.DeveloperRewards.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.CommunityPool.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.RelayerIncentives.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.DistributionProportions.Size()
	n += 1 + l + sovParams(uint64(l))
	if len(m.WeightedDeveloperRewardsReceivers) > 0 {
		for _, e := range m.WeightedDeveloperRewardsReceivers {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.SupplementAmount) > 0 {
		for _, e := range m.SupplementAmount {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WeightedAddress) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: WeightedAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WeightedAddress: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Weight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *DistributionProportions) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: DistributionProportions: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DistributionProportions: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeveloperRewards", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DeveloperRewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommunityPool", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CommunityPool.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RelayerIncentives", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RelayerIncentives.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributionProportions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DistributionProportions.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WeightedDeveloperRewardsReceivers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WeightedDeveloperRewardsReceivers = append(m.WeightedDeveloperRewardsReceivers, WeightedAddress{})
			if err := m.WeightedDeveloperRewardsReceivers[len(m.WeightedDeveloperRewardsReceivers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SupplementAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SupplementAmount = append(m.SupplementAmount, types.Coin{})
			if err := m.SupplementAmount[len(m.SupplementAmount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
