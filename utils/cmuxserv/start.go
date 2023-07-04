package cmuxserv

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type HttpServerConfig struct {
	Addr        string
	HttpHandler http.Handler
}

type RpcServerConfig struct {
	RpcOptions     []grpc.ServerOption
	RpcReadTimeOut time.Duration
	RpcRegister    func(*grpc.Server)
}

func (r *RpcServerConfig) refresh() {
	if r == nil {
		return
	}

	rpcTimeOut := time.Second * 2
	if r.RpcReadTimeOut > 0 {
		rpcTimeOut = r.RpcReadTimeOut
	}

	r.RpcReadTimeOut = rpcTimeOut
}

type CmuxConfig struct {
	RpcServerConfig  *RpcServerConfig
	HttpServerConfig *HttpServerConfig
	WsServerConfig   *HttpServerConfig

	ShutdownTime time.Duration
	StartTime    time.Duration
}

func (c *CmuxConfig) refresh() {
	c.RpcServerConfig.refresh()

	shutdownTime := time.Second * 2
	startTime := time.Millisecond * 200

	if c.StartTime > 0 {
		startTime = c.StartTime
	}
	if c.ShutdownTime > 0 {
		shutdownTime = c.ShutdownTime
	}

	c.ShutdownTime = shutdownTime
	c.StartTime = startTime
}

func newCmuxServer(c *CmuxConfig) *CmuxServer {
	addr := c.HttpServerConfig.Addr
	handler := c.HttpServerConfig.HttpHandler
	shutdownTime := c.ShutdownTime
	startTime := c.StartTime
	isGraceful := false

	baseName := filepath.Base(os.Args[0])
	GracefulEnv = baseName + "_" + GracefulEnv
	CountKeyEnv = baseName + "_" + CountKeyEnv

	for _, v := range os.Environ() {
		if v == GracefulEnv {
			isGraceful = true
			break
		}
	}

	serv := &CmuxServer{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},

		isGraceful:   isGraceful,
		shutdownTime: shutdownTime,
		startTime:    startTime,
	}
	if c.RpcServerConfig != nil {
		rpcOptions := c.RpcServerConfig.RpcOptions
		serv.rpcReadTimeOut = c.RpcServerConfig.RpcReadTimeOut
		serv.grpcServer = grpc.NewServer(rpcOptions...)
	}

	if c.WsServerConfig != nil {
		wsConfig := c.WsServerConfig
		serv.wsServer = &http.Server{
			Handler: wsConfig.HttpHandler,
		}
	}

	serv.inherit()
	return serv
}

func handleSignal(serv *CmuxServer) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for s := range ch {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			serv.Stop()
			return
		case syscall.SIGHUP:
			serv.GracefulStop()
			return
		default:
			return
		}
	}
}

func StartCmuxServer(c *CmuxConfig) *CmuxServer {
	c.refresh()
	if c.HttpServerConfig == nil {
		panic("error: HttpServerConfig is nil")
	}

	serv := newCmuxServer(c)

	// Create the main listener.
	l, err := serv.listen(c.HttpServerConfig.Addr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	m := cmux.New(l)

	// 目前遇到一个问题，当在测试环境基于 cmonitor -grestart 服务时，会停在grpc.GracefulStop
	// 本地无法复现；单独开启一个服务发布到测试环境也无法复现
	// 出现这个问题时是rpc服务没有接收过请求，一旦接收到一次请求就没问题了。（如果服务基于rpc接口探活也不会遇到这种问题）
	// 猜测是一直处于读等待，这里加个超时
	if serv.grpcServer != nil {
		rpcTimeOut := serv.rpcReadTimeOut
		m.SetReadTimeout(rpcTimeOut)

		// Match connections in order:
		// First grpc, then HTTP, and otherwise Go RPC/TCP.
		serv.grpcL = m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

		// Register gprc methods
		c.RpcServerConfig.RpcRegister(serv.grpcServer)
	}

	if serv.wsServer != nil {
		serv.websocketL = m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))
	}
	serv.httpL = m.Match(cmux.HTTP1Fast())

	serv.m = m

	go func() {
		if err := serv.Serve(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()

	handleSignal(serv)

	return serv
}
