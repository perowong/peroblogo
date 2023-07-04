package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/perowong/peroblogo/internal/http/controller"
	"github.com/perowong/peroblogo/internal/model"
)

func AuthToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		ctxUtils := &controller.GinCtxUtils{Context: ctx}
		if token == "" {
			ctxUtils.ReplyFail(controller.ErrCodeToken, "bad token")
			ctx.Abort()
			return
		}

		daoObj := model.NewModel()
		userToken, err := daoObj.GetUserToken(token)
		if err != nil {
			ctxUtils.ReplyFail(controller.ErrCodeToken, "bad token")
			ctx.Abort()
			return
		}

		if time.Now().After(userToken.ExpireTime) {
			ctxUtils.ReplyFail(controller.ErrCodeToken, "token is expired")
			ctx.Abort()
			return
		}

		go func() {
			_ = daoObj.UpdateUserTokenExpireTime(token)
		}()

		ctx.Next()
	}
}
