package middleware

import (
	"context"
	"grpc-exmpl/internal/service"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthMiddleware struct {
	userService service.UserService
}

func NewAuthMiddleware(userService service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

// UnaryInterceptor for unary calls
func (a *AuthMiddleware) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Skip auth for certain methods
	if a.isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	// Extract token from metadata
	token, err := a.extractToken(ctx)
	if err != nil {
		return nil, err
	}

	// Validate token
	claims, err := a.userService.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Add user info to context
	ctx = a.addUserToContext(ctx, claims.UserID, claims.Username, claims.Email)

	return handler(ctx, req)
}

// StreamInterceptor for streaming calls
func (a *AuthMiddleware) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// Skip auth for certain methods
	if a.isPublicMethod(info.FullMethod) {
		return handler(srv, ss)
	}

	// Extract token from metadata
	token, err := a.extractToken(ss.Context())
	if err != nil {
		return err
	}

	// Validate token
	claims, err := a.userService.ValidateToken(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// Add user info to context
	ctx := a.addUserToContext(ss.Context(), claims.UserID, claims.Username, claims.Email)
	wrappedStream := &wrappedServerStream{ss, ctx}

	return handler(srv, wrappedStream)
}

// isPublicMethod checks if the method doesn't require authentication
func (a *AuthMiddleware) isPublicMethod(method string) bool {
	publicMethods := []string{
		"/user.UserService/Register",
		"/user.UserService/Login",
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}
	return false
}

// extractToken extracts JWT token from gRPC metadata
func (a *AuthMiddleware) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	// Extract token from "Bearer <token>" format
	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	return parts[1], nil
}

// addUserToContext adds user information to context
func (a *AuthMiddleware) addUserToContext(ctx context.Context, userID int64, username, email string) context.Context {
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "username", username)
	ctx = context.WithValue(ctx, "email", email)
	return ctx
}

// wrappedServerStream wraps grpc.ServerStream with custom context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// Helper functions to extract user info from context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value("user_id").(int64)
	return userID, ok
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value("username").(string)
	return username, ok
}

func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)
	return email, ok
}
