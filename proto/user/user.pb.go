package user

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterRequest represents registration request parameters.
type RegisterRequest struct {
	Username string
	Email    string
	Password string
	FullName string
}

// RegisterResponse represents registration response.
type RegisterResponse struct {
	Success bool
	Message string
	User    *UserData
}

// LoginRequest represents login request parameters.
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represents login response data.
type LoginResponse struct {
	Success bool
	Message string
	Token   string
	User    *UserData
}

// GetProfileRequest is used to request a user profile by token.
type GetProfileRequest struct {
	Token string
}

// GetProfileResponse returns a user profile.
type GetProfileResponse struct {
	Success bool
	Message string
	User    *UserData
}

// UserData describes a user entity.
type UserData struct {
	Id        int64
	Username  string
	Email     string
	FullName  string
	CreatedAt string
	UpdatedAt string
}

// UserServiceClient defines the gRPC client API for UserService service.
type UserServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error)
}

type userServiceClient struct{ cc grpc.ClientConnInterface }

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/user.UserService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/user.UserService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error) {
	out := new(GetProfileResponse)
	err := c.cc.Invoke(ctx, "/user.UserService/GetProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer defines the gRPC server API for UserService service.
type UserServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
}

// UnimplementedUserServiceServer can be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserServiceServer) GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/GetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetProfile(ctx, req.(*GetProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Register", Handler: _UserService_Register_Handler},
		{MethodName: "Login", Handler: _UserService_Login_Handler},
		{MethodName: "GetProfile", Handler: _UserService_GetProfile_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/user/user.proto",
}
