// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.15.8
// source: product.proto

package src

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

const (
	ProductServices_Add_FullMethodName    = "/cost_control.ProductServices/Add"
	ProductServices_Update_FullMethodName = "/cost_control.ProductServices/Update"
	ProductServices_Get_FullMethodName    = "/cost_control.ProductServices/Get"
	ProductServices_Search_FullMethodName = "/cost_control.ProductServices/Search"
	ProductServices_Delete_FullMethodName = "/cost_control.ProductServices/Delete"
	ProductServices_Report_FullMethodName = "/cost_control.ProductServices/Report"
)

// ProductServicesClient is the client API for ProductServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductServicesClient interface {
	Add(ctx context.Context, in *CreateProductRequest, opts ...grpc.CallOption) (*CreateProductResponse, error)
	Update(ctx context.Context, in *UpdateProductRequest, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*GetProductResponse, error)
	Search(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*SearchProductResponse, error)
	Delete(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*Response, error)
	Report(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ReportResponse, error)
}

type productServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServicesClient(cc grpc.ClientConnInterface) ProductServicesClient {
	return &productServicesClient{cc}
}

func (c *productServicesClient) Add(ctx context.Context, in *CreateProductRequest, opts ...grpc.CallOption) (*CreateProductResponse, error) {
	out := new(CreateProductResponse)
	err := c.cc.Invoke(ctx, ProductServices_Add_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServicesClient) Update(ctx context.Context, in *UpdateProductRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, ProductServices_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServicesClient) Get(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*GetProductResponse, error) {
	out := new(GetProductResponse)
	err := c.cc.Invoke(ctx, ProductServices_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServicesClient) Search(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*SearchProductResponse, error) {
	out := new(SearchProductResponse)
	err := c.cc.Invoke(ctx, ProductServices_Search_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServicesClient) Delete(ctx context.Context, in *ProductRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, ProductServices_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServicesClient) Report(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ReportResponse, error) {
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, ProductServices_Report_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductServicesServer is the server API for ProductServices service.
// All implementations must embed UnimplementedProductServicesServer
// for forward compatibility
type ProductServicesServer interface {
	Add(context.Context, *CreateProductRequest) (*CreateProductResponse, error)
	Update(context.Context, *UpdateProductRequest) (*Response, error)
	Get(context.Context, *ProductRequest) (*GetProductResponse, error)
	Search(context.Context, *Filter) (*SearchProductResponse, error)
	Delete(context.Context, *ProductRequest) (*Response, error)
	Report(context.Context, *Filter) (*ReportResponse, error)
	mustEmbedUnimplementedProductServicesServer()
}

// UnimplementedProductServicesServer must be embedded to have forward compatible implementations.
type UnimplementedProductServicesServer struct {
}

func (UnimplementedProductServicesServer) Add(context.Context, *CreateProductRequest) (*CreateProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedProductServicesServer) Update(context.Context, *UpdateProductRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedProductServicesServer) Get(context.Context, *ProductRequest) (*GetProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedProductServicesServer) Search(context.Context, *Filter) (*SearchProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedProductServicesServer) Delete(context.Context, *ProductRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedProductServicesServer) Report(context.Context, *Filter) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
func (UnimplementedProductServicesServer) mustEmbedUnimplementedProductServicesServer() {}

// UnsafeProductServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductServicesServer will
// result in compilation errors.
type UnsafeProductServicesServer interface {
	mustEmbedUnimplementedProductServicesServer()
}

func RegisterProductServicesServer(s grpc.ServiceRegistrar, srv ProductServicesServer) {
	s.RegisterService(&ProductServices_ServiceDesc, srv)
}

func _ProductServices_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServicesServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductServices_Add_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServicesServer).Add(ctx, req.(*CreateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductServices_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServicesServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductServices_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServicesServer).Update(ctx, req.(*UpdateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductServices_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServicesServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductServices_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServicesServer).Get(ctx, req.(*ProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductServices_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServicesServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductServices_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServicesServer).Search(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductServices_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServicesServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductServices_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServicesServer).Delete(ctx, req.(*ProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductServices_Report_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServicesServer).Report(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductServices_Report_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServicesServer).Report(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductServices_ServiceDesc is the grpc.ServiceDesc for ProductServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cost_control.ProductServices",
	HandlerType: (*ProductServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _ProductServices_Add_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ProductServices_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ProductServices_Get_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _ProductServices_Search_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ProductServices_Delete_Handler,
		},
		{
			MethodName: "Report",
			Handler:    _ProductServices_Report_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "product.proto",
}
