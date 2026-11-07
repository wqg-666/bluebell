package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//1.注册业务路由
	//从路由走到controller这一层，controller负责处理路由，参数校验，请求转发
	r.POST("/signup", controller.SignUpHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	return r
}
