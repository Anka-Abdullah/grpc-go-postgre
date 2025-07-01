package grpc

import (
	"context"
	"grpc-exmpl/internal/model"
	"grpc-exmpl/internal/service"
	pb "grpc-exmpl/proto/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Convert proto request to model
	registerReq := &model.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	// Call service
	user, err := h.userService.Register(registerReq)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
			User:    nil,
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert model to proto response
	userProto := &pb.UserData{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		User:    userProto,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Convert proto request to model
	loginReq := &model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service
	loginResp, err := h.userService.Login(loginReq)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
			Token:   "",
			User:    nil,
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	// Convert model to proto response
	userProto := &pb.UserData{
		Id:        loginResp.User.ID,
		Username:  loginResp.User.Username,
		Email:     loginResp.User.Email,
		FullName:  loginResp.User.FullName,
		CreatedAt: loginResp.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: loginResp.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   loginResp.Token,
		User:    userProto,
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	// Call service
	user, err := h.userService.GetProfile(req.Token)
	if err != nil {
		return &pb.GetProfileResponse{
			Success: false,
			Message: err.Error(),
			User:    nil,
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	// Convert model to proto response
	userProto := &pb.UserData{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return &pb.GetProfileResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		User:    userProto,
	}, nil
}
