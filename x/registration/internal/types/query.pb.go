// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: registration/v1beta1/query.proto

package types

import (
	bytes "bytes"
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type QueryEncryptedSeedRequest struct {
	PubKey []byte `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
}

func (m *QueryEncryptedSeedRequest) Reset()         { *m = QueryEncryptedSeedRequest{} }
func (m *QueryEncryptedSeedRequest) String() string { return proto.CompactTextString(m) }
func (*QueryEncryptedSeedRequest) ProtoMessage()    {}
func (*QueryEncryptedSeedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc9f9bba42351659, []int{0}
}
func (m *QueryEncryptedSeedRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEncryptedSeedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEncryptedSeedRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEncryptedSeedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEncryptedSeedRequest.Merge(m, src)
}
func (m *QueryEncryptedSeedRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryEncryptedSeedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEncryptedSeedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEncryptedSeedRequest proto.InternalMessageInfo

type QueryEncryptedSeedResponse struct {
	EncryptedSeed []byte `protobuf:"bytes,1,opt,name=encrypted_seed,json=encryptedSeed,proto3" json:"encrypted_seed,omitempty"`
}

func (m *QueryEncryptedSeedResponse) Reset()         { *m = QueryEncryptedSeedResponse{} }
func (m *QueryEncryptedSeedResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEncryptedSeedResponse) ProtoMessage()    {}
func (*QueryEncryptedSeedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc9f9bba42351659, []int{1}
}
func (m *QueryEncryptedSeedResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEncryptedSeedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEncryptedSeedResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEncryptedSeedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEncryptedSeedResponse.Merge(m, src)
}
func (m *QueryEncryptedSeedResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryEncryptedSeedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEncryptedSeedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEncryptedSeedResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryEncryptedSeedRequest)(nil), "trst.x.registration.v1beta1.QueryEncryptedSeedRequest")
	proto.RegisterType((*QueryEncryptedSeedResponse)(nil), "trst.x.registration.v1beta1.QueryEncryptedSeedResponse")
}

func init() { proto.RegisterFile("registration/v1beta1/query.proto", fileDescriptor_cc9f9bba42351659) }

var fileDescriptor_cc9f9bba42351659 = []byte{
	// 407 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x51, 0xcd, 0xaa, 0xd3, 0x40,
	0x18, 0x4d, 0x94, 0x56, 0x18, 0xac, 0xc2, 0x20, 0xfe, 0xa4, 0x65, 0x28, 0x45, 0xa5, 0x9b, 0xce,
	0x58, 0x95, 0xba, 0x57, 0xba, 0xea, 0xca, 0xea, 0x42, 0xdc, 0x94, 0xa4, 0xfd, 0x1c, 0x43, 0xdb,
	0xcc, 0x34, 0x33, 0x91, 0x84, 0xe2, 0xc6, 0x27, 0x10, 0x7c, 0x89, 0x3e, 0x82, 0x8f, 0xd0, 0x65,
	0xc1, 0x8d, 0xcb, 0x7b, 0xd3, 0xfb, 0x20, 0x97, 0x4c, 0xd2, 0x4b, 0x0a, 0xb9, 0x85, 0xcb, 0xdd,
	0x4d, 0x72, 0xbe, 0xf3, 0x9d, 0x73, 0xbe, 0x83, 0xda, 0x21, 0x70, 0x5f, 0xe9, 0xd0, 0xd5, 0xbe,
	0x08, 0xd8, 0x8f, 0xbe, 0x07, 0xda, 0xed, 0xb3, 0x55, 0x04, 0x61, 0x42, 0x65, 0x28, 0xb4, 0xc0,
	0x4d, 0x1d, 0x2a, 0x4d, 0x63, 0x5a, 0x1e, 0xa4, 0xc5, 0xa0, 0xf3, 0x88, 0x0b, 0x2e, 0xcc, 0x1c,
	0xcb, 0x5e, 0x39, 0xc5, 0x69, 0x72, 0x21, 0xf8, 0x02, 0x98, 0xf9, 0xf2, 0xa2, 0x6f, 0x0c, 0x96,
	0x52, 0x17, 0xfb, 0x9c, 0x56, 0x01, 0xba, 0xd2, 0x67, 0x6e, 0x10, 0x08, 0x6d, 0x36, 0xaa, 0x02,
	0x25, 0x95, 0x7e, 0x96, 0x8a, 0xe7, 0x78, 0xe7, 0x2d, 0x7a, 0xf6, 0x31, 0x33, 0x37, 0x0c, 0xa6,
	0x61, 0x22, 0x35, 0xcc, 0x3e, 0x01, 0xcc, 0xc6, 0xb0, 0x8a, 0x40, 0x69, 0xfc, 0x04, 0xdd, 0x93,
	0x91, 0x37, 0x99, 0x43, 0xf2, 0xd4, 0x6e, 0xdb, 0xdd, 0xfb, 0xe3, 0xba, 0x8c, 0xbc, 0x11, 0x24,
	0x9d, 0x0f, 0xc8, 0xa9, 0x62, 0x29, 0x29, 0x02, 0x05, 0xf8, 0x05, 0x7a, 0x00, 0x07, 0x60, 0xa2,
	0x00, 0x66, 0x05, 0xbb, 0x01, 0xe5, 0xf1, 0xd7, 0x9b, 0xbb, 0xa8, 0x66, 0xb6, 0x60, 0x8e, 0x6a,
	0x9f, 0xe3, 0x11, 0x24, 0xf8, 0x31, 0xcd, 0xc3, 0xd0, 0x43, 0x52, 0x3a, 0xcc, 0x92, 0x3a, 0x6d,
	0x7a, 0xe2, 0x68, 0x34, 0x73, 0xf4, 0xfc, 0xd7, 0xbf, 0x8b, 0x3f, 0x77, 0x08, 0x6e, 0xb1, 0xca,
	0xc0, 0x3a, 0xee, 0xcd, 0x21, 0xc1, 0x6b, 0xf4, 0x70, 0x5c, 0x82, 0x6f, 0x27, 0x49, 0x8d, 0x64,
	0x17, 0xbf, 0xac, 0x96, 0x2c, 0xff, 0x34, 0xe2, 0x7f, 0x6d, 0xd4, 0x38, 0x3a, 0x18, 0x1e, 0x9c,
	0xd4, 0xb8, 0xb6, 0x17, 0xe7, 0xdd, 0x8d, 0x79, 0x79, 0x33, 0x9d, 0x81, 0xb1, 0xfc, 0x0a, 0xd3,
	0x6a, 0xcb, 0x57, 0xfd, 0xf4, 0xb2, 0xd6, 0xd8, 0xba, 0x28, 0xff, 0xe7, 0xfb, 0x2f, 0xdb, 0x73,
	0x62, 0x6d, 0x52, 0x62, 0x6f, 0x53, 0x62, 0xef, 0x52, 0x62, 0x9f, 0xa5, 0xc4, 0xfe, 0xbd, 0x27,
	0xd6, 0x6e, 0x4f, 0xac, 0xff, 0x7b, 0x62, 0x7d, 0x1d, 0x70, 0x5f, 0x7f, 0x8f, 0x3c, 0x3a, 0x15,
	0x4b, 0x96, 0x99, 0x5b, 0xb8, 0x9e, 0x32, 0x0f, 0x16, 0x1f, 0x6b, 0xf9, 0x81, 0x86, 0x30, 0x70,
	0x17, 0x4c, 0x27, 0x12, 0x94, 0x57, 0x37, 0x67, 0x7f, 0x73, 0x19, 0x00, 0x00, 0xff, 0xff, 0x61,
	0xb9, 0x1b, 0xe3, 0x38, 0x03, 0x00, 0x00,
}

func (this *QueryEncryptedSeedRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*QueryEncryptedSeedRequest)
	if !ok {
		that2, ok := that.(QueryEncryptedSeedRequest)
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
	if !bytes.Equal(this.PubKey, that1.PubKey) {
		return false
	}
	return true
}
func (this *QueryEncryptedSeedResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*QueryEncryptedSeedResponse)
	if !ok {
		that2, ok := that.(QueryEncryptedSeedResponse)
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
	if !bytes.Equal(this.EncryptedSeed, that1.EncryptedSeed) {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Returns the key used for transactions
	TxKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Key, error)
	// Returns the key used for registration
	RegistrationKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Key, error)
	// Returns the encrypted seed for a registered node by public key
	EncryptedSeed(ctx context.Context, in *QueryEncryptedSeedRequest, opts ...grpc.CallOption) (*QueryEncryptedSeedResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) TxKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Key, error) {
	out := new(Key)
	err := c.cc.Invoke(ctx, "/trst.x.registration.v1beta1.Query/TxKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) RegistrationKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Key, error) {
	out := new(Key)
	err := c.cc.Invoke(ctx, "/trst.x.registration.v1beta1.Query/RegistrationKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EncryptedSeed(ctx context.Context, in *QueryEncryptedSeedRequest, opts ...grpc.CallOption) (*QueryEncryptedSeedResponse, error) {
	out := new(QueryEncryptedSeedResponse)
	err := c.cc.Invoke(ctx, "/trst.x.registration.v1beta1.Query/EncryptedSeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Returns the key used for transactions
	TxKey(context.Context, *emptypb.Empty) (*Key, error)
	// Returns the key used for registration
	RegistrationKey(context.Context, *emptypb.Empty) (*Key, error)
	// Returns the encrypted seed for a registered node by public key
	EncryptedSeed(context.Context, *QueryEncryptedSeedRequest) (*QueryEncryptedSeedResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) TxKey(ctx context.Context, req *emptypb.Empty) (*Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TxKey not implemented")
}
func (*UnimplementedQueryServer) RegistrationKey(ctx context.Context, req *emptypb.Empty) (*Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegistrationKey not implemented")
}
func (*UnimplementedQueryServer) EncryptedSeed(ctx context.Context, req *QueryEncryptedSeedRequest) (*QueryEncryptedSeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EncryptedSeed not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_TxKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TxKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trst.x.registration.v1beta1.Query/TxKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TxKey(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_RegistrationKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RegistrationKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trst.x.registration.v1beta1.Query/RegistrationKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RegistrationKey(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EncryptedSeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEncryptedSeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EncryptedSeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trst.x.registration.v1beta1.Query/EncryptedSeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EncryptedSeed(ctx, req.(*QueryEncryptedSeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "trst.x.registration.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TxKey",
			Handler:    _Query_TxKey_Handler,
		},
		{
			MethodName: "RegistrationKey",
			Handler:    _Query_RegistrationKey_Handler,
		},
		{
			MethodName: "EncryptedSeed",
			Handler:    _Query_EncryptedSeed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registration/v1beta1/query.proto",
}

func (m *QueryEncryptedSeedRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEncryptedSeedRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEncryptedSeedRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryEncryptedSeedResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEncryptedSeedResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEncryptedSeedResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EncryptedSeed) > 0 {
		i -= len(m.EncryptedSeed)
		copy(dAtA[i:], m.EncryptedSeed)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.EncryptedSeed)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryEncryptedSeedRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryEncryptedSeedResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EncryptedSeed)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryEncryptedSeedRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryEncryptedSeedRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEncryptedSeedRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryEncryptedSeedResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryEncryptedSeedResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEncryptedSeedResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EncryptedSeed", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EncryptedSeed = append(m.EncryptedSeed[:0], dAtA[iNdEx:postIndex]...)
			if m.EncryptedSeed == nil {
				m.EncryptedSeed = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
