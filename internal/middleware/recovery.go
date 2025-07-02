package middleware

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryMiddleware recovers from panics in gRPC handlers.
type RecoveryMiddleware struct {
	logger *logrus.Entry
}

// NewRecoveryMiddleware creates a new RecoveryMiddleware.
func NewRecoveryMiddleware(logger *logrus.Entry) *RecoveryMiddleware {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &RecoveryMiddleware{logger: logger}
}

// UnaryInterceptor recovers from panics in unary RPCs.
func (m *RecoveryMiddleware) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			m.logger.WithFields(logrus.Fields{
				"method": info.FullMethod,
				"panic":  r,
			}).Error("recovered from panic in unary handler")
			err = status.Error(codes.Internal, "internal server error")
		}
	}()

	return handler(ctx, req)
}

// StreamInterceptor recovers from panics in streaming RPCs.
func (m *RecoveryMiddleware) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			m.logger.WithFields(logrus.Fields{
				"method": info.FullMethod,
				"panic":  r,
			}).Error("recovered from panic in stream handler")
			err = status.Error(codes.Internal, "internal server error")
		}
	}()

	return handler(srv, ss)
}
