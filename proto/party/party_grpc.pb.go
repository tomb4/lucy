// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package partyProto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PartyServiceClient is the client API for PartyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PartyServiceClient interface {
	// 加入派对
	JoinMetaParty(ctx context.Context, in *JoinMetaPartyReq, opts ...grpc.CallOption) (*MetaParty, error)
	// 退出派对
	ExitMetaParty(ctx context.Context, in *ExitMetaPartyReq, opts ...grpc.CallOption) (*Nil, error)
}

type partyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPartyServiceClient(cc grpc.ClientConnInterface) PartyServiceClient {
	return &partyServiceClient{cc}
}

func (c *partyServiceClient) JoinMetaParty(ctx context.Context, in *JoinMetaPartyReq, opts ...grpc.CallOption) (*MetaParty, error) {
	out := new(MetaParty)
	err := c.cc.Invoke(ctx, "/partyProto.PartyService/JoinMetaParty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partyServiceClient) ExitMetaParty(ctx context.Context, in *ExitMetaPartyReq, opts ...grpc.CallOption) (*Nil, error) {
	out := new(Nil)
	err := c.cc.Invoke(ctx, "/partyProto.PartyService/ExitMetaParty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PartyServiceServer is the server API for PartyService service.
// All implementations must embed UnimplementedPartyServiceServer
// for forward compatibility
type PartyServiceServer interface {
	// 加入派对
	JoinMetaParty(context.Context, *JoinMetaPartyReq) (*MetaParty, error)
	// 退出派对
	ExitMetaParty(context.Context, *ExitMetaPartyReq) (*Nil, error)
	mustEmbedUnimplementedPartyServiceServer()
}

// UnimplementedPartyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPartyServiceServer struct {
}

func (UnimplementedPartyServiceServer) JoinMetaParty(context.Context, *JoinMetaPartyReq) (*MetaParty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinMetaParty not implemented")
}
func (UnimplementedPartyServiceServer) ExitMetaParty(context.Context, *ExitMetaPartyReq) (*Nil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExitMetaParty not implemented")
}
func (UnimplementedPartyServiceServer) mustEmbedUnimplementedPartyServiceServer() {}

// UnsafePartyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PartyServiceServer will
// result in compilation errors.
type UnsafePartyServiceServer interface {
	mustEmbedUnimplementedPartyServiceServer()
}

func RegisterPartyServiceServer(s grpc.ServiceRegistrar, srv PartyServiceServer) {
	s.RegisterService(&PartyService_ServiceDesc, srv)
}

func _PartyService_JoinMetaParty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinMetaPartyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyServiceServer).JoinMetaParty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/partyProto.PartyService/JoinMetaParty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyServiceServer).JoinMetaParty(ctx, req.(*JoinMetaPartyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PartyService_ExitMetaParty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExitMetaPartyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartyServiceServer).ExitMetaParty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/partyProto.PartyService/ExitMetaParty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartyServiceServer).ExitMetaParty(ctx, req.(*ExitMetaPartyReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PartyService_ServiceDesc is the grpc.ServiceDesc for PartyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PartyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "partyProto.PartyService",
	HandlerType: (*PartyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JoinMetaParty",
			Handler:    _PartyService_JoinMetaParty_Handler,
		},
		{
			MethodName: "ExitMetaParty",
			Handler:    _PartyService_ExitMetaParty_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "party.proto",
}
