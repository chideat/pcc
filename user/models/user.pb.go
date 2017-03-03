// Code generated by protoc-gen-go.
// source: user.proto
// DO NOT EDIT!

/*
Package models is a generated protocol buffer package.

It is generated from these files:
	user.proto

It has these top-level messages:
	User
*/
package models

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	Id          uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Password    string `protobuf:"bytes,6,opt,name=password" json:"password,omitempty"`
	Token       string `protobuf:"bytes,7,opt,name=token" json:"token,omitempty"`
	CreatedUtc  int64  `protobuf:"varint,14,opt,name=created_utc,json=createdUtc" json:"created_utc,omitempty"`
	ModifiedUtc int64  `protobuf:"varint,15,opt,name=modified_utc,json=modifiedUtc" json:"modified_utc,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *User) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *User) GetCreatedUtc() int64 {
	if m != nil {
		return m.CreatedUtc
	}
	return 0
}

func (m *User) GetModifiedUtc() int64 {
	if m != nil {
		return m.ModifiedUtc
	}
	return 0
}

func init() {
	proto.RegisterType((*User)(nil), "models.User")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UserRPC service

type UserRPCClient interface {
	GetUserById(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
}

type userRPCClient struct {
	cc *grpc.ClientConn
}

func NewUserRPCClient(cc *grpc.ClientConn) UserRPCClient {
	return &userRPCClient{cc}
}

func (c *userRPCClient) GetUserById(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/models.UserRPC/GetUserById", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserRPC service

type UserRPCServer interface {
	GetUserById(context.Context, *User) (*User, error)
}

func RegisterUserRPCServer(s *grpc.Server, srv UserRPCServer) {
	s.RegisterService(&_UserRPC_serviceDesc, srv)
}

func _UserRPC_GetUserById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRPCServer).GetUserById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/models.UserRPC/GetUserById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRPCServer).GetUserById(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "models.UserRPC",
	HandlerType: (*UserRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserById",
			Handler:    _UserRPC_GetUserById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}

func init() { proto.RegisterFile("user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0xcd, 0x4f, 0x49, 0xcd, 0x29, 0x56, 0x5a,
	0xc0, 0xc8, 0xc5, 0x12, 0x5a, 0x9c, 0x5a, 0x24, 0xc4, 0xc7, 0xc5, 0x94, 0x99, 0x22, 0xc1, 0xa8,
	0xc0, 0xa8, 0xc1, 0x12, 0xc4, 0x94, 0x99, 0x22, 0x24, 0xc4, 0xc5, 0x92, 0x97, 0x98, 0x9b, 0x2a,
	0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66, 0x0b, 0x49, 0x71, 0x71, 0x14, 0x24, 0x16, 0x17,
	0x97, 0xe7, 0x17, 0xa5, 0x48, 0xb0, 0x81, 0xc5, 0xe1, 0x7c, 0x21, 0x11, 0x2e, 0xd6, 0x92, 0xfc,
	0xec, 0xd4, 0x3c, 0x09, 0x76, 0xb0, 0x04, 0x84, 0x23, 0x24, 0xcf, 0xc5, 0x9d, 0x5c, 0x94, 0x9a,
	0x58, 0x92, 0x9a, 0x12, 0x5f, 0x5a, 0x92, 0x2c, 0xc1, 0xa7, 0xc0, 0xa8, 0xc1, 0x1c, 0xc4, 0x05,
	0x15, 0x0a, 0x2d, 0x49, 0x16, 0x52, 0xe4, 0xe2, 0xc9, 0xcd, 0x4f, 0xc9, 0x4c, 0xcb, 0x84, 0xaa,
	0xe0, 0x07, 0xab, 0xe0, 0x86, 0x89, 0x85, 0x96, 0x24, 0x1b, 0x99, 0x71, 0xb1, 0x83, 0x5c, 0x18,
	0x14, 0xe0, 0x2c, 0xa4, 0xcd, 0xc5, 0xed, 0x9e, 0x5a, 0x02, 0xe2, 0x39, 0x55, 0x7a, 0xa6, 0x08,
	0xf1, 0xe8, 0x41, 0x7c, 0xa1, 0x07, 0x12, 0x91, 0x42, 0xe1, 0x29, 0x31, 0x24, 0xb1, 0x81, 0x7d,
	0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xae, 0x10, 0xe9, 0x40, 0xf7, 0x00, 0x00, 0x00,
}