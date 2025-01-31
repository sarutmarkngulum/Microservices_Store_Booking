// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: menu.proto

package services

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MenuService_CreateMenuItem_FullMethodName            = "/services.MenuService/CreateMenuItem"
	MenuService_UpdateMenuItem_FullMethodName            = "/services.MenuService/UpdateMenuItem"
	MenuService_DeleteMenuItem_FullMethodName            = "/services.MenuService/DeleteMenuItem"
	MenuService_GetMenuItems_FullMethodName              = "/services.MenuService/GetMenuItems"
	MenuService_GetMenuItemById_FullMethodName           = "/services.MenuService/GetMenuItemById"
	MenuService_CreateMenuSet_FullMethodName             = "/services.MenuService/CreateMenuSet"
	MenuService_UpdateMenuSet_FullMethodName             = "/services.MenuService/UpdateMenuSet"
	MenuService_DeleteMenuSet_FullMethodName             = "/services.MenuService/DeleteMenuSet"
	MenuService_GetMenuSets_FullMethodName               = "/services.MenuService/GetMenuSets"
	MenuService_GetMenuSetById_FullMethodName            = "/services.MenuService/GetMenuSetById"
	MenuService_CreateMenuSetItem_FullMethodName         = "/services.MenuService/CreateMenuSetItem"
	MenuService_GetMenuSetItems_FullMethodName           = "/services.MenuService/GetMenuSetItems"
	MenuService_GetMenuSetItemByMenuSetID_FullMethodName = "/services.MenuService/GetMenuSetItemByMenuSetID"
	MenuService_UpdateMenuSetItem_FullMethodName         = "/services.MenuService/UpdateMenuSetItem"
	MenuService_DeleteMenuSetItem_FullMethodName         = "/services.MenuService/DeleteMenuSetItem"
	MenuService_UploadImage_FullMethodName               = "/services.MenuService/UploadImage"
	MenuService_DeleteImage_FullMethodName               = "/services.MenuService/DeleteImage"
)

// MenuServiceClient is the client API for MenuService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MenuServiceClient interface {
	// Handle Menu
	CreateMenuItem(ctx context.Context, in *CreateMenuItemRequest, opts ...grpc.CallOption) (*CreateMenuItemResponse, error)
	UpdateMenuItem(ctx context.Context, in *UpdateMenuItemRequest, opts ...grpc.CallOption) (*UpdateMenuItemResponse, error)
	DeleteMenuItem(ctx context.Context, in *DeleteMenuItemRequest, opts ...grpc.CallOption) (*DeleteMenuItemResponse, error)
	GetMenuItems(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuItemList, error)
	GetMenuItemById(ctx context.Context, in *GetMenuItemByIdRequest, opts ...grpc.CallOption) (*MenuItem, error)
	// Handle Menu Set
	CreateMenuSet(ctx context.Context, in *CreateMenuSetRequest, opts ...grpc.CallOption) (*CreateMenuSetResponse, error)
	UpdateMenuSet(ctx context.Context, in *UpdateMenuSetRequest, opts ...grpc.CallOption) (*UpdateMenuSetResponse, error)
	DeleteMenuSet(ctx context.Context, in *DeleteMenuSetRequest, opts ...grpc.CallOption) (*DeleteMenuSetResponse, error)
	GetMenuSets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuSetList, error)
	GetMenuSetById(ctx context.Context, in *GetMenuSetByIdRequest, opts ...grpc.CallOption) (*MenuSet, error)
	// Handle Menu Set Item
	CreateMenuSetItem(ctx context.Context, in *CreateMenuSetItemRequest, opts ...grpc.CallOption) (*CreateMenuSetItemResponse, error)
	GetMenuSetItems(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuSetItemList, error)
	GetMenuSetItemByMenuSetID(ctx context.Context, in *GetMenuSetItemByIdRequest, opts ...grpc.CallOption) (*MenuSetItemList, error)
	UpdateMenuSetItem(ctx context.Context, in *UpdateMenuSetItemRequest, opts ...grpc.CallOption) (*UpdateMenuSetItemResponse, error)
	DeleteMenuSetItem(ctx context.Context, in *DeleteMenuSetItemRequest, opts ...grpc.CallOption) (*DeleteMenuSetItemResponse, error)
	// Handle Upload Img On Cloud
	UploadImage(ctx context.Context, in *UploadImageRequest, opts ...grpc.CallOption) (*UploadImageResponse, error)
	DeleteImage(ctx context.Context, in *DeleteImageRequest, opts ...grpc.CallOption) (*DeleteImageResponse, error)
}

type menuServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMenuServiceClient(cc grpc.ClientConnInterface) MenuServiceClient {
	return &menuServiceClient{cc}
}

func (c *menuServiceClient) CreateMenuItem(ctx context.Context, in *CreateMenuItemRequest, opts ...grpc.CallOption) (*CreateMenuItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMenuItemResponse)
	err := c.cc.Invoke(ctx, MenuService_CreateMenuItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) UpdateMenuItem(ctx context.Context, in *UpdateMenuItemRequest, opts ...grpc.CallOption) (*UpdateMenuItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMenuItemResponse)
	err := c.cc.Invoke(ctx, MenuService_UpdateMenuItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) DeleteMenuItem(ctx context.Context, in *DeleteMenuItemRequest, opts ...grpc.CallOption) (*DeleteMenuItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMenuItemResponse)
	err := c.cc.Invoke(ctx, MenuService_DeleteMenuItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) GetMenuItems(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuItemList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MenuItemList)
	err := c.cc.Invoke(ctx, MenuService_GetMenuItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) GetMenuItemById(ctx context.Context, in *GetMenuItemByIdRequest, opts ...grpc.CallOption) (*MenuItem, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MenuItem)
	err := c.cc.Invoke(ctx, MenuService_GetMenuItemById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) CreateMenuSet(ctx context.Context, in *CreateMenuSetRequest, opts ...grpc.CallOption) (*CreateMenuSetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMenuSetResponse)
	err := c.cc.Invoke(ctx, MenuService_CreateMenuSet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) UpdateMenuSet(ctx context.Context, in *UpdateMenuSetRequest, opts ...grpc.CallOption) (*UpdateMenuSetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMenuSetResponse)
	err := c.cc.Invoke(ctx, MenuService_UpdateMenuSet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) DeleteMenuSet(ctx context.Context, in *DeleteMenuSetRequest, opts ...grpc.CallOption) (*DeleteMenuSetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMenuSetResponse)
	err := c.cc.Invoke(ctx, MenuService_DeleteMenuSet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) GetMenuSets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuSetList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MenuSetList)
	err := c.cc.Invoke(ctx, MenuService_GetMenuSets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) GetMenuSetById(ctx context.Context, in *GetMenuSetByIdRequest, opts ...grpc.CallOption) (*MenuSet, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MenuSet)
	err := c.cc.Invoke(ctx, MenuService_GetMenuSetById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) CreateMenuSetItem(ctx context.Context, in *CreateMenuSetItemRequest, opts ...grpc.CallOption) (*CreateMenuSetItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMenuSetItemResponse)
	err := c.cc.Invoke(ctx, MenuService_CreateMenuSetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) GetMenuSetItems(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MenuSetItemList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MenuSetItemList)
	err := c.cc.Invoke(ctx, MenuService_GetMenuSetItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) GetMenuSetItemByMenuSetID(ctx context.Context, in *GetMenuSetItemByIdRequest, opts ...grpc.CallOption) (*MenuSetItemList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MenuSetItemList)
	err := c.cc.Invoke(ctx, MenuService_GetMenuSetItemByMenuSetID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) UpdateMenuSetItem(ctx context.Context, in *UpdateMenuSetItemRequest, opts ...grpc.CallOption) (*UpdateMenuSetItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMenuSetItemResponse)
	err := c.cc.Invoke(ctx, MenuService_UpdateMenuSetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) DeleteMenuSetItem(ctx context.Context, in *DeleteMenuSetItemRequest, opts ...grpc.CallOption) (*DeleteMenuSetItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMenuSetItemResponse)
	err := c.cc.Invoke(ctx, MenuService_DeleteMenuSetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) UploadImage(ctx context.Context, in *UploadImageRequest, opts ...grpc.CallOption) (*UploadImageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadImageResponse)
	err := c.cc.Invoke(ctx, MenuService_UploadImage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *menuServiceClient) DeleteImage(ctx context.Context, in *DeleteImageRequest, opts ...grpc.CallOption) (*DeleteImageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteImageResponse)
	err := c.cc.Invoke(ctx, MenuService_DeleteImage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MenuServiceServer is the server API for MenuService service.
// All implementations must embed UnimplementedMenuServiceServer
// for forward compatibility.
type MenuServiceServer interface {
	// Handle Menu
	CreateMenuItem(context.Context, *CreateMenuItemRequest) (*CreateMenuItemResponse, error)
	UpdateMenuItem(context.Context, *UpdateMenuItemRequest) (*UpdateMenuItemResponse, error)
	DeleteMenuItem(context.Context, *DeleteMenuItemRequest) (*DeleteMenuItemResponse, error)
	GetMenuItems(context.Context, *emptypb.Empty) (*MenuItemList, error)
	GetMenuItemById(context.Context, *GetMenuItemByIdRequest) (*MenuItem, error)
	// Handle Menu Set
	CreateMenuSet(context.Context, *CreateMenuSetRequest) (*CreateMenuSetResponse, error)
	UpdateMenuSet(context.Context, *UpdateMenuSetRequest) (*UpdateMenuSetResponse, error)
	DeleteMenuSet(context.Context, *DeleteMenuSetRequest) (*DeleteMenuSetResponse, error)
	GetMenuSets(context.Context, *emptypb.Empty) (*MenuSetList, error)
	GetMenuSetById(context.Context, *GetMenuSetByIdRequest) (*MenuSet, error)
	// Handle Menu Set Item
	CreateMenuSetItem(context.Context, *CreateMenuSetItemRequest) (*CreateMenuSetItemResponse, error)
	GetMenuSetItems(context.Context, *emptypb.Empty) (*MenuSetItemList, error)
	GetMenuSetItemByMenuSetID(context.Context, *GetMenuSetItemByIdRequest) (*MenuSetItemList, error)
	UpdateMenuSetItem(context.Context, *UpdateMenuSetItemRequest) (*UpdateMenuSetItemResponse, error)
	DeleteMenuSetItem(context.Context, *DeleteMenuSetItemRequest) (*DeleteMenuSetItemResponse, error)
	// Handle Upload Img On Cloud
	UploadImage(context.Context, *UploadImageRequest) (*UploadImageResponse, error)
	DeleteImage(context.Context, *DeleteImageRequest) (*DeleteImageResponse, error)
	mustEmbedUnimplementedMenuServiceServer()
}

// UnimplementedMenuServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMenuServiceServer struct{}

func (UnimplementedMenuServiceServer) CreateMenuItem(context.Context, *CreateMenuItemRequest) (*CreateMenuItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMenuItem not implemented")
}
func (UnimplementedMenuServiceServer) UpdateMenuItem(context.Context, *UpdateMenuItemRequest) (*UpdateMenuItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMenuItem not implemented")
}
func (UnimplementedMenuServiceServer) DeleteMenuItem(context.Context, *DeleteMenuItemRequest) (*DeleteMenuItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMenuItem not implemented")
}
func (UnimplementedMenuServiceServer) GetMenuItems(context.Context, *emptypb.Empty) (*MenuItemList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuItems not implemented")
}
func (UnimplementedMenuServiceServer) GetMenuItemById(context.Context, *GetMenuItemByIdRequest) (*MenuItem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuItemById not implemented")
}
func (UnimplementedMenuServiceServer) CreateMenuSet(context.Context, *CreateMenuSetRequest) (*CreateMenuSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMenuSet not implemented")
}
func (UnimplementedMenuServiceServer) UpdateMenuSet(context.Context, *UpdateMenuSetRequest) (*UpdateMenuSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMenuSet not implemented")
}
func (UnimplementedMenuServiceServer) DeleteMenuSet(context.Context, *DeleteMenuSetRequest) (*DeleteMenuSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMenuSet not implemented")
}
func (UnimplementedMenuServiceServer) GetMenuSets(context.Context, *emptypb.Empty) (*MenuSetList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuSets not implemented")
}
func (UnimplementedMenuServiceServer) GetMenuSetById(context.Context, *GetMenuSetByIdRequest) (*MenuSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuSetById not implemented")
}
func (UnimplementedMenuServiceServer) CreateMenuSetItem(context.Context, *CreateMenuSetItemRequest) (*CreateMenuSetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMenuSetItem not implemented")
}
func (UnimplementedMenuServiceServer) GetMenuSetItems(context.Context, *emptypb.Empty) (*MenuSetItemList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuSetItems not implemented")
}
func (UnimplementedMenuServiceServer) GetMenuSetItemByMenuSetID(context.Context, *GetMenuSetItemByIdRequest) (*MenuSetItemList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuSetItemByMenuSetID not implemented")
}
func (UnimplementedMenuServiceServer) UpdateMenuSetItem(context.Context, *UpdateMenuSetItemRequest) (*UpdateMenuSetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMenuSetItem not implemented")
}
func (UnimplementedMenuServiceServer) DeleteMenuSetItem(context.Context, *DeleteMenuSetItemRequest) (*DeleteMenuSetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMenuSetItem not implemented")
}
func (UnimplementedMenuServiceServer) UploadImage(context.Context, *UploadImageRequest) (*UploadImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}
func (UnimplementedMenuServiceServer) DeleteImage(context.Context, *DeleteImageRequest) (*DeleteImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteImage not implemented")
}
func (UnimplementedMenuServiceServer) mustEmbedUnimplementedMenuServiceServer() {}
func (UnimplementedMenuServiceServer) testEmbeddedByValue()                     {}

// UnsafeMenuServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MenuServiceServer will
// result in compilation errors.
type UnsafeMenuServiceServer interface {
	mustEmbedUnimplementedMenuServiceServer()
}

func RegisterMenuServiceServer(s grpc.ServiceRegistrar, srv MenuServiceServer) {
	// If the following call pancis, it indicates UnimplementedMenuServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MenuService_ServiceDesc, srv)
}

func _MenuService_CreateMenuItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMenuItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).CreateMenuItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_CreateMenuItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).CreateMenuItem(ctx, req.(*CreateMenuItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_UpdateMenuItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMenuItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).UpdateMenuItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_UpdateMenuItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).UpdateMenuItem(ctx, req.(*UpdateMenuItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_DeleteMenuItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMenuItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).DeleteMenuItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_DeleteMenuItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).DeleteMenuItem(ctx, req.(*DeleteMenuItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_GetMenuItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).GetMenuItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_GetMenuItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).GetMenuItems(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_GetMenuItemById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMenuItemByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).GetMenuItemById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_GetMenuItemById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).GetMenuItemById(ctx, req.(*GetMenuItemByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_CreateMenuSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMenuSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).CreateMenuSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_CreateMenuSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).CreateMenuSet(ctx, req.(*CreateMenuSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_UpdateMenuSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMenuSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).UpdateMenuSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_UpdateMenuSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).UpdateMenuSet(ctx, req.(*UpdateMenuSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_DeleteMenuSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMenuSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).DeleteMenuSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_DeleteMenuSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).DeleteMenuSet(ctx, req.(*DeleteMenuSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_GetMenuSets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).GetMenuSets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_GetMenuSets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).GetMenuSets(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_GetMenuSetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMenuSetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).GetMenuSetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_GetMenuSetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).GetMenuSetById(ctx, req.(*GetMenuSetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_CreateMenuSetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMenuSetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).CreateMenuSetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_CreateMenuSetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).CreateMenuSetItem(ctx, req.(*CreateMenuSetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_GetMenuSetItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).GetMenuSetItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_GetMenuSetItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).GetMenuSetItems(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_GetMenuSetItemByMenuSetID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMenuSetItemByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).GetMenuSetItemByMenuSetID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_GetMenuSetItemByMenuSetID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).GetMenuSetItemByMenuSetID(ctx, req.(*GetMenuSetItemByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_UpdateMenuSetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMenuSetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).UpdateMenuSetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_UpdateMenuSetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).UpdateMenuSetItem(ctx, req.(*UpdateMenuSetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_DeleteMenuSetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMenuSetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).DeleteMenuSetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_DeleteMenuSetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).DeleteMenuSetItem(ctx, req.(*DeleteMenuSetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_UploadImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).UploadImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_UploadImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).UploadImage(ctx, req.(*UploadImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MenuService_DeleteImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MenuServiceServer).DeleteImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MenuService_DeleteImage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MenuServiceServer).DeleteImage(ctx, req.(*DeleteImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MenuService_ServiceDesc is the grpc.ServiceDesc for MenuService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MenuService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.MenuService",
	HandlerType: (*MenuServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMenuItem",
			Handler:    _MenuService_CreateMenuItem_Handler,
		},
		{
			MethodName: "UpdateMenuItem",
			Handler:    _MenuService_UpdateMenuItem_Handler,
		},
		{
			MethodName: "DeleteMenuItem",
			Handler:    _MenuService_DeleteMenuItem_Handler,
		},
		{
			MethodName: "GetMenuItems",
			Handler:    _MenuService_GetMenuItems_Handler,
		},
		{
			MethodName: "GetMenuItemById",
			Handler:    _MenuService_GetMenuItemById_Handler,
		},
		{
			MethodName: "CreateMenuSet",
			Handler:    _MenuService_CreateMenuSet_Handler,
		},
		{
			MethodName: "UpdateMenuSet",
			Handler:    _MenuService_UpdateMenuSet_Handler,
		},
		{
			MethodName: "DeleteMenuSet",
			Handler:    _MenuService_DeleteMenuSet_Handler,
		},
		{
			MethodName: "GetMenuSets",
			Handler:    _MenuService_GetMenuSets_Handler,
		},
		{
			MethodName: "GetMenuSetById",
			Handler:    _MenuService_GetMenuSetById_Handler,
		},
		{
			MethodName: "CreateMenuSetItem",
			Handler:    _MenuService_CreateMenuSetItem_Handler,
		},
		{
			MethodName: "GetMenuSetItems",
			Handler:    _MenuService_GetMenuSetItems_Handler,
		},
		{
			MethodName: "GetMenuSetItemByMenuSetID",
			Handler:    _MenuService_GetMenuSetItemByMenuSetID_Handler,
		},
		{
			MethodName: "UpdateMenuSetItem",
			Handler:    _MenuService_UpdateMenuSetItem_Handler,
		},
		{
			MethodName: "DeleteMenuSetItem",
			Handler:    _MenuService_DeleteMenuSetItem_Handler,
		},
		{
			MethodName: "UploadImage",
			Handler:    _MenuService_UploadImage_Handler,
		},
		{
			MethodName: "DeleteImage",
			Handler:    _MenuService_DeleteImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "menu.proto",
}
