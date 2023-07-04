package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/internal/http/controller"
	"github.com/perowong/peroblogo/internal/http/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupHttpRouters(r *gin.Engine) {
	// groupTokenRequired := r.Use(middleware.AuthToken())
	{
		r.Any("health", controller.Health)
	}

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
