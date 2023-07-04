package grpcregister

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// 参考文档：https://github.com/grpc/grpc/blob/master/doc/health-checking.md

// HealthImpl 健康检查实现
type HealthImpl struct{}

// Watch
// A client can call the Watch method to perform a streaming health-check.
// The server will immediately send back a message indicating the current serving status.
// It will then subsequently send a new message whenever the service's serving status changes.
func (h *HealthImpl) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return nil
}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
// A client can query the server’s health status by calling the Check method, and a deadline should be set on the rpc.
// The client can optionally set the service name it wants to query for health status.
// The suggested format of service name is package_names.ServiceName, such as grpc.health.v1.Health.
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	resp := grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}
	return &resp, nil
}

func RegisterHealthService(s *grpc.Server) {
	//grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})
}
