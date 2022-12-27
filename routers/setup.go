package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/conf"
	"github.com/perowong/peroblogo/controller"
)

func SetupRouters() *gin.Engine {
	if conf.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	groupRouters(r)

	return r
}

func groupRouters(r *gin.Engine) {
	group := r.Group("comment")
	{
		group.POST("/add", controller.AddComment)
		group.POST("/list", controller.ListComment)
	}
}
