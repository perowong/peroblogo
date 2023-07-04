package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/perowong/peroblogo/utils/consul"
)

type EnvType string

const (
	Development EnvType = "dev"
	Test        EnvType = "test"
	Production  EnvType = "prod"
)

type GRpcLogger struct {
	Logger *zap.Logger
	Opts   []grpc_zap.Option
}

type gRpcConnOptions struct {
	env        EnvType
	devAddr    string
	retryCodes []codes.Code
	logger     *GRpcLogger
}

type GRpcConnOption interface {
	apply(*gRpcConnOptions)
}

type funcGRpcOption struct {
	f func(*gRpcConnOptions)
}

func (fdo *funcGRpcOption) apply(do *gRpcConnOptions) {
	fdo.f(do)
}

func newFuncGRpcOption(f func(*gRpcConnOptions)) *funcGRpcOption {
	return &funcGRpcOption{
		f: f,
	}
}

func WithDevelopment(env EnvType, devAddr string) GRpcConnOption {
	return newFuncGRpcOption(func(o *gRpcConnOptions) {
		o.env = env
		o.devAddr = devAddr
	})
}

func WithRetryCodes(codes ...codes.Code) GRpcConnOption {
	return newFuncGRpcOption(func(o *gRpcConnOptions) {
		o.retryCodes = append(o.retryCodes, codes...)
	})
}

func WithLogger(logger *GRpcLogger) GRpcConnOption {
	return newFuncGRpcOption(func(o *gRpcConnOptions) {
		o.logger = logger
	})
}

func NewGRpcConn(serviceName string, gRpcOption ...GRpcConnOption) (conn *grpc.ClientConn, err error) {
	gRpcOpts := &gRpcConnOptions{}
	for _, opt := range gRpcOption {
		opt.apply(gRpcOpts)
	}

	opts := []retry.CallOption{
		retry.WithBackoff(retry.BackoffLinear(10 * time.Millisecond)),
		retry.WithMax(5),
	}

	retryCodes := []codes.Code{codes.NotFound, codes.Aborted, codes.Unavailable}
	if len(gRpcOpts.retryCodes) > 0 {
		retryCodes = append(retryCodes, gRpcOpts.retryCodes...)
	}
	opts = append(opts, retry.WithCodes(retryCodes...))

	backoffConfig := backoff.DefaultConfig
	backoffConfig.MaxDelay = 5 * time.Second
	conParams := grpc.ConnectParams{
		Backoff:           backoffConfig,
		MinConnectTimeout: 1 * time.Second,
	}

	chainStreamClient := getChainStreamClient(opts, gRpcOpts.logger)
	chainUnaryClient := getChainUnaryClient(opts, gRpcOpts.logger)

	addr := getDialAddr(serviceName, gRpcOpts.devAddr, gRpcOpts.env)
	conn, err = grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpcMiddleware.ChainStreamClient(chainStreamClient...)),
		grpc.WithUnaryInterceptor(grpcMiddleware.ChainUnaryClient(chainUnaryClient...)),
		grpc.WithConnectParams(conParams),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	)

	return
}

func getDialAddr(consulName, devAddr string, env EnvType) (target string) {
	ip := "127.0.0.1"

	k8sNodeIP := os.Getenv("K8S_NODE_IP")
	k8sNodeIP = strings.TrimSpace(k8sNodeIP)

	if k8sNodeIP != "" {
		ip = k8sNodeIP
	}
	if env == Development && len(devAddr) > 0 {
		return devAddr
	} else {
		return "consul://" + ip + ":8500/" + consulName
	}
}

func getChainStreamClient(opts []retry.CallOption, logger *GRpcLogger) []grpc.StreamClientInterceptor {
	clArr := make([]grpc.StreamClientInterceptor, 0)
	clArr = append(clArr, grpcPrometheus.StreamClientInterceptor)
	clArr = append(clArr, retry.StreamClientInterceptor(opts...))
	if logger != nil {
		clArr = append(clArr, grpc_zap.StreamClientInterceptor(logger.Logger, logger.Opts...))
	}

	return clArr
}

func getChainUnaryClient(opts []retry.CallOption, logger *GRpcLogger) []grpc.UnaryClientInterceptor {
	clArr := make([]grpc.UnaryClientInterceptor, 0)
	clArr = append(clArr, grpcPrometheus.UnaryClientInterceptor)
	clArr = append(clArr, retry.UnaryClientInterceptor(opts...))
	// clArr = append(clArr, clientInterceptor())
	if logger != nil {
		clArr = append(clArr, grpc_zap.UnaryClientInterceptor(logger.Logger, logger.Opts...))
	}

	return clArr
}
