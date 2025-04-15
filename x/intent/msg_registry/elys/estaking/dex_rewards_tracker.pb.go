// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/estaking/dex_rewards_tracker.proto

package estaking

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// DexRewardsTracker is used for tracking rewards for stakers and LPs, all
// amount here is in USDC
type DexRewardsTracker struct {
	// Number of blocks since start of epoch (distribution epoch)
	NumBlocks int64 `protobuf:"varint,1,opt,name=num_blocks,json=numBlocks,proto3" json:"num_blocks,omitempty"`
	// Accumulated amount at distribution epoch - recalculated at every
	// distribution epoch
	Amount cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=amount,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"amount"`
}

func (m *DexRewardsTracker) Reset()         { *m = DexRewardsTracker{} }
func (m *DexRewardsTracker) String() string { return proto.CompactTextString(m) }
func (*DexRewardsTracker) ProtoMessage()    {}
func (*DexRewardsTracker) Descriptor() ([]byte, []int) {
	return fileDescriptor_061875a2058b444a, []int{0}
}
func (m *DexRewardsTracker) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DexRewardsTracker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DexRewardsTracker.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DexRewardsTracker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DexRewardsTracker.Merge(m, src)
}
func (m *DexRewardsTracker) XXX_Size() int {
	return m.Size()
}
func (m *DexRewardsTracker) XXX_DiscardUnknown() {
	xxx_messageInfo_DexRewardsTracker.DiscardUnknown(m)
}

var xxx_messageInfo_DexRewardsTracker proto.InternalMessageInfo

func (m *DexRewardsTracker) GetNumBlocks() int64 {
	if m != nil {
		return m.NumBlocks
	}
	return 0
}

func init() {
	proto.RegisterType((*DexRewardsTracker)(nil), "elys.estaking.DexRewardsTracker")
}

func init() {
	proto.RegisterFile("elys/estaking/dex_rewards_tracker.proto", fileDescriptor_061875a2058b444a)
}

var fileDescriptor_061875a2058b444a = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x4f, 0x4a, 0xc3, 0x40,
	0x14, 0xc6, 0x33, 0x0a, 0x85, 0x0e, 0xb8, 0xb0, 0xb8, 0xa8, 0x15, 0xa7, 0xc5, 0x8d, 0xdd, 0x98,
	0x41, 0x3c, 0x80, 0x50, 0xb2, 0x11, 0x5c, 0x05, 0x57, 0x82, 0x84, 0xc9, 0x64, 0x98, 0x86, 0x64,
	0x66, 0x64, 0xe6, 0x05, 0x93, 0x85, 0x77, 0xf0, 0x18, 0x1e, 0xc0, 0x43, 0x74, 0x59, 0x5c, 0x89,
	0x8b, 0x22, 0xc9, 0x45, 0xa4, 0x99, 0xb8, 0x70, 0xf7, 0xbd, 0xef, 0x7d, 0xfc, 0xde, 0x1f, 0x7c,
	0x29, 0xca, 0xc6, 0x51, 0xe1, 0x80, 0x15, 0xb9, 0x96, 0x34, 0x13, 0x75, 0x62, 0xc5, 0x0b, 0xb3,
	0x99, 0x4b, 0xc0, 0x32, 0x5e, 0x08, 0x1b, 0x3e, 0x5b, 0x03, 0x66, 0x72, 0xb4, 0x0f, 0x86, 0x7f,
	0xc1, 0xd9, 0x29, 0x37, 0x4e, 0x19, 0x97, 0xf4, 0x4d, 0xea, 0x0b, 0x9f, 0x9c, 0x9d, 0x48, 0x23,
	0x8d, 0xf7, 0xf7, 0xca, 0xbb, 0x17, 0xaf, 0xf8, 0x38, 0x12, 0x75, 0xec, 0xd9, 0x0f, 0x1e, 0x3d,
	0x39, 0xc7, 0x58, 0x57, 0x2a, 0x49, 0x4b, 0xc3, 0x0b, 0x37, 0x45, 0x0b, 0xb4, 0x3c, 0x8c, 0xc7,
	0xba, 0x52, 0xab, 0xde, 0x98, 0xdc, 0xe1, 0x11, 0x53, 0xa6, 0xd2, 0x30, 0x3d, 0x58, 0xa0, 0xe5,
	0x78, 0x75, 0xbd, 0xd9, 0xcd, 0x83, 0xef, 0xdd, 0xfc, 0xcc, 0xcf, 0x73, 0x59, 0x11, 0xe6, 0x86,
	0x2a, 0x06, 0xeb, 0xf0, 0x5e, 0x48, 0xc6, 0x9b, 0x48, 0xf0, 0xcf, 0x8f, 0x2b, 0x3c, 0xac, 0x13,
	0x09, 0x1e, 0x0f, 0x80, 0xd5, 0xd3, 0x7b, 0x4b, 0xd0, 0xa6, 0x25, 0x68, 0xdb, 0x12, 0xf4, 0xd3,
	0x12, 0xf4, 0xd6, 0x91, 0x60, 0xdb, 0x91, 0xe0, 0xab, 0x23, 0xc1, 0xe3, 0xad, 0xcc, 0x61, 0x5d,
	0xa5, 0x21, 0x37, 0x8a, 0x82, 0x75, 0x50, 0xb2, 0xd4, 0xd1, 0x5c, 0x83, 0xd0, 0x60, 0x68, 0x3d,
	0x28, 0xaa, 0x9c, 0x4c, 0xac, 0x90, 0xb9, 0x03, 0xdb, 0xd0, 0x7f, 0x7f, 0x4b, 0x47, 0xfd, 0x91,
	0x37, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x01, 0xde, 0x48, 0xae, 0x4f, 0x01, 0x00, 0x00,
}

func (this *DexRewardsTracker) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DexRewardsTracker)
	if !ok {
		that2, ok := that.(DexRewardsTracker)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.NumBlocks != that1.NumBlocks {
		return false
	}
	if !this.Amount.Equal(that1.Amount) {
		return false
	}
	return true
}
func (m *DexRewardsTracker) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DexRewardsTracker) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DexRewardsTracker) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDexRewardsTracker(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.NumBlocks != 0 {
		i = encodeVarintDexRewardsTracker(dAtA, i, uint64(m.NumBlocks))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintDexRewardsTracker(dAtA []byte, offset int, v uint64) int {
	offset -= sovDexRewardsTracker(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DexRewardsTracker) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.NumBlocks != 0 {
		n += 1 + sovDexRewardsTracker(uint64(m.NumBlocks))
	}
	l = m.Amount.Size()
	n += 1 + l + sovDexRewardsTracker(uint64(l))
	return n
}

func sovDexRewardsTracker(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDexRewardsTracker(x uint64) (n int) {
	return sovDexRewardsTracker(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DexRewardsTracker) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDexRewardsTracker
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
			return fmt.Errorf("proto: DexRewardsTracker: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DexRewardsTracker: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumBlocks", wireType)
			}
			m.NumBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDexRewardsTracker
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumBlocks |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDexRewardsTracker
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
				return ErrInvalidLengthDexRewardsTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDexRewardsTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDexRewardsTracker(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDexRewardsTracker
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
func skipDexRewardsTracker(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDexRewardsTracker
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
					return 0, ErrIntOverflowDexRewardsTracker
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
					return 0, ErrIntOverflowDexRewardsTracker
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
				return 0, ErrInvalidLengthDexRewardsTracker
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDexRewardsTracker
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDexRewardsTracker
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDexRewardsTracker        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDexRewardsTracker          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDexRewardsTracker = fmt.Errorf("proto: unexpected end of group")
)
