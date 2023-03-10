package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/conf"
	"github.com/perowong/peroblogo/controller"
	_ "github.com/perowong/peroblogo/docs"
	"github.com/perowong/peroblogo/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouters() *gin.Engine {
	if conf.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	groupRouters(r)

	return r
}

func groupRouters(r *gin.Engine) {
	group := r.Group("api").Use(middleware.AuthToken())
	groupWithoutAuth := r.Group("api")
	{
		group.POST("comment/add", controller.AddComment)
	}
	{
		groupWithoutAuth.POST("comment/list", controller.ListComment)
		groupWithoutAuth.POST("user/login", controller.Login)
	}

	r.GET("/api/swagger-doc/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
