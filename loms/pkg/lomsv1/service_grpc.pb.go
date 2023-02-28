// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: service.proto

package lomsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LomsV1Client is the client API for LomsV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LomsV1Client interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderID, error)
	ListOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*ListOrderResponse, error)
	OrderPayed(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CancelOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Stocks(ctx context.Context, in *SKU, opts ...grpc.CallOption) (*StocksList, error)
}

type lomsV1Client struct {
	cc grpc.ClientConnInterface
}

func NewLomsV1Client(cc grpc.ClientConnInterface) LomsV1Client {
	return &lomsV1Client{cc}
}

func (c *lomsV1Client) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*OrderID, error) {
	out := new(OrderID)
	err := c.cc.Invoke(ctx, "/loms_v1.LomsV1/createOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsV1Client) ListOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*ListOrderResponse, error) {
	out := new(ListOrderResponse)
	err := c.cc.Invoke(ctx, "/loms_v1.LomsV1/listOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsV1Client) OrderPayed(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loms_v1.LomsV1/orderPayed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsV1Client) CancelOrder(ctx context.Context, in *OrderID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loms_v1.LomsV1/cancelOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsV1Client) Stocks(ctx context.Context, in *SKU, opts ...grpc.CallOption) (*StocksList, error) {
	out := new(StocksList)
	err := c.cc.Invoke(ctx, "/loms_v1.LomsV1/stocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LomsV1Server is the server API for LomsV1 service.
// All implementations must embed UnimplementedLomsV1Server
// for forward compatibility
type LomsV1Server interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*OrderID, error)
	ListOrder(context.Context, *OrderID) (*ListOrderResponse, error)
	OrderPayed(context.Context, *OrderID) (*emptypb.Empty, error)
	CancelOrder(context.Context, *OrderID) (*emptypb.Empty, error)
	Stocks(context.Context, *SKU) (*StocksList, error)
	mustEmbedUnimplementedLomsV1Server()
}

// UnimplementedLomsV1Server must be embedded to have forward compatible implementations.
type UnimplementedLomsV1Server struct {
}

func (UnimplementedLomsV1Server) CreateOrder(context.Context, *CreateOrderRequest) (*OrderID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedLomsV1Server) ListOrder(context.Context, *OrderID) (*ListOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrder not implemented")
}
func (UnimplementedLomsV1Server) OrderPayed(context.Context, *OrderID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderPayed not implemented")
}
func (UnimplementedLomsV1Server) CancelOrder(context.Context, *OrderID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedLomsV1Server) Stocks(context.Context, *SKU) (*StocksList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stocks not implemented")
}
func (UnimplementedLomsV1Server) mustEmbedUnimplementedLomsV1Server() {}

// UnsafeLomsV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LomsV1Server will
// result in compilation errors.
type UnsafeLomsV1Server interface {
	mustEmbedUnimplementedLomsV1Server()
}

func RegisterLomsV1Server(s grpc.ServiceRegistrar, srv LomsV1Server) {
	s.RegisterService(&LomsV1_ServiceDesc, srv)
}

func _LomsV1_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsV1Server).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LomsV1/createOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsV1Server).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LomsV1_ListOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsV1Server).ListOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LomsV1/listOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsV1Server).ListOrder(ctx, req.(*OrderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LomsV1_OrderPayed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsV1Server).OrderPayed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LomsV1/orderPayed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsV1Server).OrderPayed(ctx, req.(*OrderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LomsV1_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsV1Server).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LomsV1/cancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsV1Server).CancelOrder(ctx, req.(*OrderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LomsV1_Stocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SKU)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsV1Server).Stocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms_v1.LomsV1/stocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsV1Server).Stocks(ctx, req.(*SKU))
	}
	return interceptor(ctx, in, info, handler)
}

// LomsV1_ServiceDesc is the grpc.ServiceDesc for LomsV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LomsV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loms_v1.LomsV1",
	HandlerType: (*LomsV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createOrder",
			Handler:    _LomsV1_CreateOrder_Handler,
		},
		{
			MethodName: "listOrder",
			Handler:    _LomsV1_ListOrder_Handler,
		},
		{
			MethodName: "orderPayed",
			Handler:    _LomsV1_OrderPayed_Handler,
		},
		{
			MethodName: "cancelOrder",
			Handler:    _LomsV1_CancelOrder_Handler,
		},
		{
			MethodName: "stocks",
			Handler:    _LomsV1_Stocks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
