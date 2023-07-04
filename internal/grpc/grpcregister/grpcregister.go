package grpcregister

import (
	"google.golang.org/grpc"
)

func Register(s *grpc.Server) {
	RegisterHealthService(s)
}
