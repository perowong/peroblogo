package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/perowong/peroblogo/conf"
	"github.com/perowong/peroblogo/internal/dao"
	"github.com/perowong/peroblogo/internal/grpc/grpcregister"
	"github.com/perowong/peroblogo/internal/http/routers"
	"github.com/perowong/peroblogo/utils/cmuxserv"
	"google.golang.org/grpc"
)

// @title Peroblogo Api doc
// @version 1.0
// @contact.name Pero Wong
// @contact.url https://i.overio.space
// @contact.email ynwangpeng@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host https://i.overio.space
// @BasePath /api

func main() {
	var env string
	flag.StringVar(&env, "env", "dev", "set the server's env, like -env=dev or -env=prod")
	flag.Parse()
	conf.InitConf(conf.EnvType(env))

	// dao.InitMysql()
	db := dao.ConnectMysql()
	defer db.Close()

	if conf.Env != conf.Development {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	routers.SetupHttpRouters(g)

	shutdownTime := time.Second * 2
	grpcReadTimeout := time.Second * 2
	if conf.Env != conf.Production {
		shutdownTime = time.Second * 5
		grpcReadTimeout = time.Millisecond * 300
	}

	rpcOptions := []grpc.ServerOption{
		grpcMiddleware.WithUnaryServerChain(),
	}
	mConfig := &cmuxserv.CmuxConfig{
		ShutdownTime: shutdownTime,
		RpcServerConfig: &cmuxserv.RpcServerConfig{
			RpcOptions:     rpcOptions,
			RpcReadTimeOut: grpcReadTimeout,
			RpcRegister:    grpcregister.Register,
		},
		HttpServerConfig: &cmuxserv.HttpServerConfig{
			Addr:        conf.C.App.Addr,
			HttpHandler: g,
		},
	}

	cmuxserv.StartCmuxServer(mConfig)
}
