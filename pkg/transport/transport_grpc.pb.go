// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: transport.proto

package transport

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

// AliasServiceClient is the client API for AliasService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AliasServiceClient interface {
	CreateAlias(ctx context.Context, in *URLRequest, opts ...grpc.CallOption) (*AliasResp, error)
	GetOrigURL(ctx context.Context, in *AliasReq, opts ...grpc.CallOption) (*OriginalURLResp, error)
}

type aliasServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAliasServiceClient(cc grpc.ClientConnInterface) AliasServiceClient {
	return &aliasServiceClient{cc}
}

func (c *aliasServiceClient) CreateAlias(ctx context.Context, in *URLRequest, opts ...grpc.CallOption) (*AliasResp, error) {
	out := new(AliasResp)
	err := c.cc.Invoke(ctx, "/transport.AliasService/CreateAlias", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aliasServiceClient) GetOrigURL(ctx context.Context, in *AliasReq, opts ...grpc.CallOption) (*OriginalURLResp, error) {
	out := new(OriginalURLResp)
	err := c.cc.Invoke(ctx, "/transport.AliasService/GetOrigURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AliasServiceServer is the server API for AliasService service.
// All implementations must embed UnimplementedAliasServiceServer
// for forward compatibility
type AliasServiceServer interface {
	CreateAlias(context.Context, *URLRequest) (*AliasResp, error)
	GetOrigURL(context.Context, *AliasReq) (*OriginalURLResp, error)
	mustEmbedUnimplementedAliasServiceServer()
}

// UnimplementedAliasServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAliasServiceServer struct {
}

func (UnimplementedAliasServiceServer) CreateAlias(context.Context, *URLRequest) (*AliasResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlias not implemented")
}
func (UnimplementedAliasServiceServer) GetOrigURL(context.Context, *AliasReq) (*OriginalURLResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrigURL not implemented")
}
func (UnimplementedAliasServiceServer) mustEmbedUnimplementedAliasServiceServer() {}

// UnsafeAliasServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AliasServiceServer will
// result in compilation errors.
type UnsafeAliasServiceServer interface {
	mustEmbedUnimplementedAliasServiceServer()
}

func RegisterAliasServiceServer(s grpc.ServiceRegistrar, srv AliasServiceServer) {
	s.RegisterService(&AliasService_ServiceDesc, srv)
}

func _AliasService_CreateAlias_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(URLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AliasServiceServer).CreateAlias(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transport.AliasService/CreateAlias",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AliasServiceServer).CreateAlias(ctx, req.(*URLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AliasService_GetOrigURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AliasReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AliasServiceServer).GetOrigURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transport.AliasService/GetOrigURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AliasServiceServer).GetOrigURL(ctx, req.(*AliasReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AliasService_ServiceDesc is the grpc.ServiceDesc for AliasService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AliasService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transport.AliasService",
	HandlerType: (*AliasServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAlias",
			Handler:    _AliasService_CreateAlias_Handler,
		},
		{
			MethodName: "GetOrigURL",
			Handler:    _AliasService_GetOrigURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transport.proto",
}
