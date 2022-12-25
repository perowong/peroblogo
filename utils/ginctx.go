package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/perowong/peroblogo/utils/errcode"
)

type GinCtxUtils struct {
	*gin.Context
}

func (c *GinCtxUtils) GetFieldValidatedErr(err error) (msg string) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok || (errs != nil && len(errs) == 0) {
		msg = err.Error()
		return
	}

	for _, e := range errs {
		msg = msg + "invalid filed " + e.Field()
		return
	}

	return
}

func (c *GinCtxUtils) GetReqObject(req interface{}) (ok bool) {
	err := c.ShouldBindBodyWith(req, binding.JSON)
	if err != nil {
		errMsg := c.GetFieldValidatedErr(err)
		c.ReplyFail(errcode.ErrCodeParam, errMsg)
		return false
	}

	return true
}

type Resp struct {
	ErrCode errcode.ErrCodeType
	ErrMsg  string
	Data    interface{}
}

func (c *GinCtxUtils) Reply(errCode errcode.ErrCodeType, errMsg string, data interface{}) {
	resp := &Resp{
		ErrCode: errCode,
		ErrMsg:  errMsg,
		Data:    data,
	}

	c.JSON(http.StatusOK, resp)
}

func (c *GinCtxUtils) ReplyOk(data interface{}) {
	c.Reply(errcode.ErrCodeOk, "ok", data)
}

func (c *GinCtxUtils) ReplyFail(errCode errcode.ErrCodeType, errMsg string) {
	c.Reply(errCode, fmt.Sprintf("fail: %s", errMsg), nil)
}

func (c *GinCtxUtils) ReplyFailParam() {
	c.ReplyFail(errcode.ErrCodeParam, errcode.CodeMap[errcode.ErrCodeParam])
}

func (c *GinCtxUtils) ReplyFailServer() {
	c.ReplyFail(errcode.ErrCodeServer, errcode.CodeMap[errcode.ErrCodeServer])
}
