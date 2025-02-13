// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: shop/v1/shop.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Shop_Auth_FullMethodName = "/shop.v1.Shop/Auth"
)

// ShopClient is the client API for Shop service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The Shop service definition.
type ShopClient interface {
	//	rpc Info(InfoRequest) returns (InfoResponse) {
	//	  option (google.api.http) = {
	//	    get: "/api/info"
	//	  };
	//	}
	//
	//	rpc SendCoin(SentTransaction) returns (SuccessResponse) {
	//	  option (google.api.http) = {
	//	    get: "/api/sendCoin"
	//	  };
	//	}
	//
	//	rpc BuyItem(Item) returns (SuccessResponse) {
	//	  option (google.api.http) = {
	//	    get: "/api/buy/{name}"
	//	  };
	//	}
	Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error)
}

type shopClient struct {
	cc grpc.ClientConnInterface
}

func NewShopClient(cc grpc.ClientConnInterface) ShopClient {
	return &shopClient{cc}
}

func (c *shopClient) Auth(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, Shop_Auth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShopServer is the server API for Shop service.
// All implementations must embed UnimplementedShopServer
// for forward compatibility.
//
// The Shop service definition.
type ShopServer interface {
	//	rpc Info(InfoRequest) returns (InfoResponse) {
	//	  option (google.api.http) = {
	//	    get: "/api/info"
	//	  };
	//	}
	//
	//	rpc SendCoin(SentTransaction) returns (SuccessResponse) {
	//	  option (google.api.http) = {
	//	    get: "/api/sendCoin"
	//	  };
	//	}
	//
	//	rpc BuyItem(Item) returns (SuccessResponse) {
	//	  option (google.api.http) = {
	//	    get: "/api/buy/{name}"
	//	  };
	//	}
	Auth(context.Context, *AuthRequest) (*AuthResponse, error)
	mustEmbedUnimplementedShopServer()
}

// UnimplementedShopServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedShopServer struct{}

func (UnimplementedShopServer) Auth(context.Context, *AuthRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}
func (UnimplementedShopServer) mustEmbedUnimplementedShopServer() {}
func (UnimplementedShopServer) testEmbeddedByValue()              {}

// UnsafeShopServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShopServer will
// result in compilation errors.
type UnsafeShopServer interface {
	mustEmbedUnimplementedShopServer()
}

func RegisterShopServer(s grpc.ServiceRegistrar, srv ShopServer) {
	// If the following call pancis, it indicates UnimplementedShopServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Shop_ServiceDesc, srv)
}

func _Shop_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShopServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shop_Auth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShopServer).Auth(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shop_ServiceDesc is the grpc.ServiceDesc for Shop service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shop_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shop.v1.Shop",
	HandlerType: (*ShopServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auth",
			Handler:    _Shop_Auth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shop/v1/shop.proto",
}
