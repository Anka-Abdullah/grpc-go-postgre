package integration

import (
	"testing"
	"time"

	apigrpc "grpc-exmpl/api/grpc"
	"grpc-exmpl/internal/service"
)

func TestGRPCServerStartStop(t *testing.T) {
	userSvc := service.NewUserService(nil, "secret")
	prodSvc := service.NewProductService(nil)

	srv := apigrpc.NewServer(userSvc, prodSvc, "0")

	done := make(chan struct{})
	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("server start: %v", err)
		}
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	srv.Stop()
	<-done
}
