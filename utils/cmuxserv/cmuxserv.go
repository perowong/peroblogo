package cmuxserv

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

const (
	// Used to indicate a graceful restart in the new process.
	envCountKey     = "LISTEN_FDS"
	fdStart     int = 3
)

var GracefulEnv = "GRACEFUL=true"
var CountKeyEnv = envCountKey

type CmuxServer struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	wsServer   *http.Server

	isGraceful     bool
	shutdownTime   time.Duration
	startTime      time.Duration
	rpcReadTimeOut time.Duration

	inherited []net.Listener
	active    []net.Listener

	grpcL      net.Listener
	httpL      net.Listener
	websocketL net.Listener
	m          cmux.CMux
}

func (srv *CmuxServer) inherit() {
	countStr := os.Getenv(CountKeyEnv)
	if countStr == "" {
		return
	}
	count, countErr := strconv.Atoi(countStr)
	if countErr != nil {
		log.Printf("found invalid count value: %s=%s", CountKeyEnv, countStr)
		return
	}

	log.Printf("found inherit fd count(%d)", count)

	for i := fdStart; i < fdStart+count; i++ {
		file := os.NewFile(uintptr(i), "listener")
		l, err := net.FileListener(file)
		if err != nil {
			_ = file.Close()
			log.Printf("error inheriting socket fd %d: %s", i, err)
			return
		}
		if err := file.Close(); err != nil {
			log.Printf("error closing inherited socket fd %d: %s", i, err)
			return
		}
		srv.inherited = append(srv.inherited, l)
	}
}

// Listen 监听GRpc端口
func (srv *CmuxServer) listen(addr string) (net.Listener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	if srv.isGraceful {
		// look for an inherited listener
		for i, l := range srv.inherited {
			if l == nil { // we nil used inherited listeners
				continue
			}
			if isSameAddr(l.Addr(), laddr) {
				log.Printf("found inherit addr(%s)", addr)
				srv.inherited[i] = nil
				srv.active = append(srv.active, l)
				return l.(*net.TCPListener), nil
			}
		}
	}
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}
	log.Printf("Listening at %s", laddr)
	srv.active = append(srv.active, l)
	return l, nil
}

func (srv *CmuxServer) fork() (err error) {
	log.Printf("grace restart...")

	// Extract the fds from the listeners.
	files := make([]*os.File, len(srv.active))
	for i, l := range srv.active {
		files[i], err = l.(*net.TCPListener).File()
		if err != nil {
			return err
		}
		defer files[i].Close()
	}

	// Pass on the environment and replace the old count key with the new one.
	var env []string
	for _, v := range os.Environ() {
		if v != GracefulEnv && !strings.HasPrefix(v, CountKeyEnv) {
			env = append(env, v)
		}
	}
	env = append(env, GracefulEnv)
	env = append(env, fmt.Sprintf("%s=%d", CountKeyEnv, len(files)))

	log.Printf("be inherited fd count(%d)", len(files))

	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	cmd.ExtraFiles = files

	log.Printf("cmd: %+v", cmd)

	err = cmd.Start()
	if err != nil {
		return
	}

	return
}

func (srv *CmuxServer) Stop() {
	log.Printf("server begin stop")
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), srv.shutdownTime)
	defer cancel()

	err := srv.httpServer.Shutdown(ctx)
	log.Printf("server finish stop http, err: %v", err)

	if srv.wsServer != nil {
		err = srv.wsServer.Shutdown(ctx)
	}
	log.Printf("server finish stop ws, err: %v", err)

	if srv.grpcServer != nil {
		srv.grpcServer.Stop()
		log.Printf("server finish stop grpc, err: %v", err)
	}

	log.Printf("server finish stop all, time past: %dns", time.Since(start).Nanoseconds())
}

func (srv *CmuxServer) GracefulStop() {
	start := time.Now()
	log.Printf("GracefulStop begin")

	err := srv.fork()
	if err != nil {
		log.Printf("start new process failed, please retry: %v\n", err)
		return
	}

	// 等待fork 进程启动
	time.Sleep(srv.startTime)

	log.Printf("GracefulStop begin http")
	ctx, cancel := context.WithTimeout(context.Background(), srv.shutdownTime)
	err = srv.httpServer.Shutdown(ctx)
	defer cancel()

	if srv.wsServer != nil {
		err = srv.wsServer.Shutdown(ctx)
	}
	log.Printf("GracefulStop end http, err: %v", err)

	if srv.grpcServer == nil {
		return
	}

	// 增加一个超时，如果未能关闭则直接停止服务
	stopped := make(chan struct{})
	go func() {
		log.Printf("GracefulStop begin grpc async")
		srv.grpcServer.GracefulStop()
		log.Printf("GracefulStop end grpc")
		close(stopped)
	}()

	timer := time.NewTimer(srv.shutdownTime)
	select {
	case <-timer.C:
		log.Printf("GracefulStop_begin_grpc_timerC")
		srv.grpcServer.Stop()
		log.Printf("GracefulStop_end_grpc_timerC")
	case <-stopped:
		log.Printf("GracefulStop_timerC_stopped")
		timer.Stop()
	}

	log.Printf("GracefulStop_end, time past: %dms", time.Since(start).Milliseconds())
}

func (srv *CmuxServer) ServeHttp(listener net.Listener) {
	if srv.httpServer == nil {
		return
	}
	//cmux.ErrListenerClosed
	if err := srv.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Printf("ServeHttp err: %+v", err)
		panic(err)
	}
}

func (srv *CmuxServer) ServeGRpc(listener net.Listener) {
	if srv.grpcServer == nil {
		return
	}
	//cmux.ErrListenerClosed
	if err := srv.grpcServer.Serve(listener); err != nil && err != cmux.ErrListenerClosed {
		log.Printf("ServeGRpc err: %+v", err)
		panic(err)
	}
}

func (srv *CmuxServer) ServeWebsocket(listener net.Listener) {
	if srv.wsServer == nil {
		return
	}
	//cmux.ErrListenerClosed
	if err := srv.wsServer.Serve(listener); err != nil && err != cmux.ErrListenerClosed {
		log.Printf("ServeWebsocket err: %+v", err)
		panic(err)
	}
}

func (srv *CmuxServer) Serve() error {
	go srv.ServeHttp(srv.httpL)

	if srv.grpcL != nil {
		go srv.ServeGRpc(srv.grpcL)
	}

	if srv.websocketL != nil {
		go srv.ServeWebsocket(srv.websocketL)
	}

	// Start serving!
	err := srv.m.Serve()
	if err != nil {
		log.Printf("cmux service err: %v", err)
	}

	return err
}
