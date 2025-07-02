package grpc

import (
	"fmt"
	"net"

	handler "grpc-exmpl/internal/handler/grpc"
	"grpc-exmpl/internal/middleware"
	"grpc-exmpl/internal/service"
	pbproduct "grpc-exmpl/proto/product"
	pbuser "grpc-exmpl/proto/user"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer     *grpc.Server
	userService    service.UserService
	productService service.ProductService
	port           string
}

func NewServer(userService service.UserService, productService service.ProductService, port string) *Server {
	return &Server{
		userService:    userService,
		productService: productService,
		port:           port,
	}
}

func (s *Server) Start() error {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(s.userService)

	// Create logrus entry for gRPC logging
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	loggingMiddleware := middleware.NewLoggingMiddleware(logrusEntry)
	recoveryMiddleware := middleware.NewRecoveryMiddleware(logrusEntry)

	// Create gRPC server with middleware
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			loggingMiddleware.UnaryInterceptor,
			recoveryMiddleware.UnaryInterceptor,
			authMiddleware.UnaryInterceptor,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			loggingMiddleware.StreamInterceptor,
			recoveryMiddleware.StreamInterceptor,
			authMiddleware.StreamInterceptor,
		)),
	)

	// Register all services
	s.registerServices()

	// Enable reflection (for development/debugging)
	reflection.Register(s.grpcServer)

	logrus.Infof("gRPC server starting on port %s", s.port)

	// Start server
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		logrus.Info("Stopping gRPC server...")
		s.grpcServer.GracefulStop()
		logrus.Info("gRPC server stopped")
	}
}

func (s *Server) registerServices() {
	// Register User service
	userHandler := handler.NewUserHandler(s.userService)
	pbuser.RegisterUserServiceServer(s.grpcServer, userHandler)

	// Register Product service
	productHandler := handler.NewProductHandler(s.productService)
	pbproduct.RegisterProductServiceServer(s.grpcServer, productHandler)

	logrus.Info("gRPC services registered successfully")
}
