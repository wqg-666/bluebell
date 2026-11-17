package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	{
	   "code": 10001,//程序中的错误码
	   "mag":xx // 提示信息
	   "data": {}//数据
	}
*/
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	//omitempty忽略掉空的值
}

func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseSucess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSucess,
		Msg:  CodeSucess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}
