package app

import (
	"github.com/gin-gonic/gin"
	"learn/gin-blog/pkg/errors"
)

type Gin struct {
	Ctx *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.Ctx.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  errors.GetMsg(errCode),
		Data: data,
	})
	return
}
