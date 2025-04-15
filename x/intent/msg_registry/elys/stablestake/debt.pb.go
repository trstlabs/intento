// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/stablestake/debt.proto

package stablestake

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

type Debt struct {
	Address               string                `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Borrowed              cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=borrowed,proto3,customtype=cosmossdk.io/math.Int" json:"borrowed"`
	InterestPaid          cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=interest_paid,json=interestPaid,proto3,customtype=cosmossdk.io/math.Int" json:"interest_paid"`
	InterestStacked       cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=interest_stacked,json=interestStacked,proto3,customtype=cosmossdk.io/math.Int" json:"interest_stacked"`
	BorrowTime            uint64                `protobuf:"varint,5,opt,name=borrow_time,json=borrowTime,proto3" json:"borrow_time,omitempty"`
	LastInterestCalcTime  uint64                `protobuf:"varint,6,opt,name=last_interest_calc_time,json=lastInterestCalcTime,proto3" json:"last_interest_calc_time,omitempty"`
	LastInterestCalcBlock uint64                `protobuf:"varint,7,opt,name=last_interest_calc_block,json=lastInterestCalcBlock,proto3" json:"last_interest_calc_block,omitempty"`
	PoolId                uint64                `protobuf:"varint,8,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
}

func (m *Debt) Reset()         { *m = Debt{} }
func (m *Debt) String() string { return proto.CompactTextString(m) }
func (*Debt) ProtoMessage()    {}
func (*Debt) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac51f1348d3e6ded, []int{0}
}
func (m *Debt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Debt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Debt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Debt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Debt.Merge(m, src)
}
func (m *Debt) XXX_Size() int {
	return m.Size()
}
func (m *Debt) XXX_DiscardUnknown() {
	xxx_messageInfo_Debt.DiscardUnknown(m)
}

var xxx_messageInfo_Debt proto.InternalMessageInfo

func (m *Debt) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Debt) GetBorrowTime() uint64 {
	if m != nil {
		return m.BorrowTime
	}
	return 0
}

func (m *Debt) GetLastInterestCalcTime() uint64 {
	if m != nil {
		return m.LastInterestCalcTime
	}
	return 0
}

func (m *Debt) GetLastInterestCalcBlock() uint64 {
	if m != nil {
		return m.LastInterestCalcBlock
	}
	return 0
}

func (m *Debt) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func init() {
	proto.RegisterType((*Debt)(nil), "elys.stablestake.Debt")
}

func init() { proto.RegisterFile("elys/stablestake/debt.proto", fileDescriptor_ac51f1348d3e6ded) }

var fileDescriptor_ac51f1348d3e6ded = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x41, 0xcb, 0xd3, 0x30,
	0x18, 0xc7, 0x5b, 0xdf, 0xb9, 0xbd, 0x46, 0xc5, 0x11, 0x36, 0x16, 0x27, 0x74, 0xc3, 0xd3, 0x40,
	0x6c, 0x0e, 0x22, 0xde, 0xab, 0x20, 0xbd, 0x8d, 0x29, 0x1e, 0x44, 0x28, 0x49, 0x13, 0xba, 0xd0,
	0xb4, 0x19, 0xc9, 0x23, 0xba, 0x6f, 0xe1, 0x87, 0xf1, 0xea, 0x7d, 0xc7, 0xe1, 0x49, 0x3c, 0x0c,
	0xd9, 0xbe, 0x88, 0xb4, 0xd9, 0x86, 0x0c, 0x2f, 0xbb, 0xe5, 0xe9, 0xef, 0xff, 0xff, 0x3d, 0x3d,
	0x3c, 0xe8, 0x89, 0xd4, 0x6b, 0x47, 0x1d, 0x30, 0xae, 0xa5, 0x03, 0x56, 0x4a, 0x2a, 0x24, 0x87,
	0x78, 0x65, 0x0d, 0x18, 0xdc, 0x6f, 0x60, 0xfc, 0x0f, 0x1c, 0x0f, 0x0a, 0x53, 0x98, 0x16, 0xd2,
	0xe6, 0xe5, 0x73, 0xe3, 0xc7, 0xb9, 0x71, 0x95, 0x71, 0x99, 0x07, 0x7e, 0xf0, 0xe8, 0xe9, 0x8f,
	0x1b, 0xd4, 0x79, 0x23, 0x39, 0x60, 0x82, 0x7a, 0x4c, 0x08, 0x2b, 0x9d, 0x23, 0xe1, 0x34, 0x9c,
	0xdd, 0x5b, 0x9c, 0x46, 0xfc, 0x16, 0xdd, 0x72, 0x63, 0xad, 0xf9, 0x22, 0x05, 0xb9, 0xd3, 0xa0,
	0xe4, 0xd9, 0x66, 0x37, 0x09, 0x7e, 0xef, 0x26, 0x43, 0xaf, 0x72, 0xa2, 0x8c, 0x95, 0xa1, 0x15,
	0x83, 0x65, 0x9c, 0xd6, 0xf0, 0xf3, 0xfb, 0x73, 0x74, 0xdc, 0x91, 0xd6, 0xb0, 0x38, 0x97, 0xf1,
	0x1c, 0x3d, 0x54, 0x35, 0x48, 0x2b, 0x1d, 0x64, 0x2b, 0xa6, 0x04, 0xb9, 0xb9, 0xde, 0xf6, 0xe0,
	0x64, 0x98, 0x33, 0x25, 0xf0, 0x07, 0xd4, 0x3f, 0x1b, 0x1d, 0xb0, 0xbc, 0x94, 0x82, 0x74, 0xae,
	0x97, 0x3e, 0x3a, 0x49, 0xde, 0x79, 0x07, 0x9e, 0xa0, 0xfb, 0xfe, 0xaf, 0x33, 0x50, 0x95, 0x24,
	0x77, 0xa7, 0xe1, 0xac, 0xb3, 0x40, 0xfe, 0xd3, 0x7b, 0x55, 0x49, 0xfc, 0x12, 0x8d, 0x34, 0x73,
	0x90, 0x9d, 0xb7, 0xe7, 0x4c, 0xe7, 0x3e, 0xdc, 0x6d, 0xc3, 0x83, 0x06, 0xa7, 0x47, 0xfa, 0x9a,
	0xe9, 0xbc, 0xad, 0xbd, 0x42, 0xe4, 0x3f, 0x35, 0xae, 0x4d, 0x5e, 0x92, 0x5e, 0xdb, 0x1b, 0x5e,
	0xf6, 0x92, 0x06, 0xe2, 0x11, 0xea, 0xad, 0x8c, 0xd1, 0x99, 0x12, 0xe4, 0xb6, 0xcd, 0x75, 0x9b,
	0x31, 0x15, 0xc9, 0xa7, 0xcd, 0x3e, 0x0a, 0xb7, 0xfb, 0x28, 0xfc, 0xb3, 0x8f, 0xc2, 0x6f, 0x87,
	0x28, 0xd8, 0x1e, 0xa2, 0xe0, 0xd7, 0x21, 0x0a, 0x3e, 0x26, 0x85, 0x82, 0xe5, 0x67, 0x1e, 0xe7,
	0xa6, 0xa2, 0x60, 0x1d, 0x68, 0xc6, 0x1d, 0x6d, 0x16, 0xd7, 0x60, 0xe8, 0xd7, 0xe3, 0x8b, 0x56,
	0xae, 0xc8, 0xac, 0x2c, 0x94, 0x03, 0xbb, 0xa6, 0x97, 0xb7, 0xc6, 0xbb, 0xed, 0x91, 0xbc, 0xf8,
	0x1b, 0x00, 0x00, 0xff, 0xff, 0x63, 0x36, 0x0c, 0xdc, 0x86, 0x02, 0x00, 0x00,
}

func (m *Debt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Debt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Debt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PoolId != 0 {
		i = encodeVarintDebt(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x40
	}
	if m.LastInterestCalcBlock != 0 {
		i = encodeVarintDebt(dAtA, i, uint64(m.LastInterestCalcBlock))
		i--
		dAtA[i] = 0x38
	}
	if m.LastInterestCalcTime != 0 {
		i = encodeVarintDebt(dAtA, i, uint64(m.LastInterestCalcTime))
		i--
		dAtA[i] = 0x30
	}
	if m.BorrowTime != 0 {
		i = encodeVarintDebt(dAtA, i, uint64(m.BorrowTime))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.InterestStacked.Size()
		i -= size
		if _, err := m.InterestStacked.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDebt(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.InterestPaid.Size()
		i -= size
		if _, err := m.InterestPaid.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDebt(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Borrowed.Size()
		i -= size
		if _, err := m.Borrowed.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDebt(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintDebt(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDebt(dAtA []byte, offset int, v uint64) int {
	offset -= sovDebt(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Debt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovDebt(uint64(l))
	}
	l = m.Borrowed.Size()
	n += 1 + l + sovDebt(uint64(l))
	l = m.InterestPaid.Size()
	n += 1 + l + sovDebt(uint64(l))
	l = m.InterestStacked.Size()
	n += 1 + l + sovDebt(uint64(l))
	if m.BorrowTime != 0 {
		n += 1 + sovDebt(uint64(m.BorrowTime))
	}
	if m.LastInterestCalcTime != 0 {
		n += 1 + sovDebt(uint64(m.LastInterestCalcTime))
	}
	if m.LastInterestCalcBlock != 0 {
		n += 1 + sovDebt(uint64(m.LastInterestCalcBlock))
	}
	if m.PoolId != 0 {
		n += 1 + sovDebt(uint64(m.PoolId))
	}
	return n
}

func sovDebt(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDebt(x uint64) (n int) {
	return sovDebt(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Debt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDebt
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
			return fmt.Errorf("proto: Debt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Debt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
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
				return ErrInvalidLengthDebt
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDebt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Borrowed", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
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
				return ErrInvalidLengthDebt
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDebt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Borrowed.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestPaid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
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
				return ErrInvalidLengthDebt
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDebt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InterestPaid.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestStacked", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
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
				return ErrInvalidLengthDebt
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDebt
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InterestStacked.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowTime", wireType)
			}
			m.BorrowTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BorrowTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastInterestCalcTime", wireType)
			}
			m.LastInterestCalcTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastInterestCalcTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastInterestCalcBlock", wireType)
			}
			m.LastInterestCalcBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastInterestCalcBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDebt
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDebt(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDebt
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
func skipDebt(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDebt
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
					return 0, ErrIntOverflowDebt
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
					return 0, ErrIntOverflowDebt
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
				return 0, ErrInvalidLengthDebt
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDebt
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDebt
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDebt        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDebt          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDebt = fmt.Errorf("proto: unexpected end of group")
)
