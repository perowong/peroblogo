package routers

import (
	"github.com/gin-gonic/gin"
)

func Setup(gin *gin.Engine) {
	group := gin.Group("comment")
	{
		group.POST("/add", AddComment)
		group.POST("/list", ListComment)
	}
}
