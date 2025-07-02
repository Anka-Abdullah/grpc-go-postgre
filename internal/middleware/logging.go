package middleware

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// LoggingMiddleware provides logging interceptors for gRPC.
type LoggingMiddleware struct {
	logger *logrus.Entry
}

// NewLoggingMiddleware creates a new LoggingMiddleware.
func NewLoggingMiddleware(logger *logrus.Entry) *LoggingMiddleware {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &LoggingMiddleware{logger: logger}
}

// UnaryInterceptor logs information about unary RPCs.
func (m *LoggingMiddleware) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	entry := m.logger.WithFields(logrus.Fields{
		"method":   info.FullMethod,
		"duration": duration.String(),
	})
	if err != nil {
		entry = entry.WithField("error", err)
		entry.Error("gRPC unary request completed with error")
	} else {
		entry.Info("gRPC unary request completed")
	}
	return resp, err
}

// StreamInterceptor logs information about streaming RPCs.
func (m *LoggingMiddleware) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	err := handler(srv, ss)
	duration := time.Since(start)

	entry := m.logger.WithFields(logrus.Fields{
		"method":   info.FullMethod,
		"duration": duration.String(),
	})
	if err != nil {
		entry = entry.WithField("error", err)
		entry.Error("gRPC stream request completed with error")
	} else {
		entry.Info("gRPC stream request completed")
	}
	return err
}
