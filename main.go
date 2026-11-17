//它的作用是告诉 Swagger：项目叫什么、后端地址是什么、接口统一前缀是什么，以及接口如何进行身份认证。
//title	文档标题	显示 Bluebell API
//version	文档版本	当前版本为 1.0
//description	文档描述	介绍这是 Bluebell 社区服务 API
//host 127.0.0.1:8080 指定服务器地址。
//schemes http 指定使用 HTTP 协议。
////Host: 127.0.0.1:8080
//BasePath API 基础路径
//securityDefinitions.apikey,定义一种 API Key 类型的认证方案，名字叫 BearerAuth。
//in header表示 Token 放在 HTTP 请求头中。
//name Authorization表示请求头中存放 Token 的字段名称是：Authorization
////Authorization: Bearer xxxxxxxxx

// @title Bluebell API
// @version 1.0
// @description Bluebell 社区服务 API 文档
// @host 127.0.0.1:8080
// @schemes http
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 Bearer 加空格和登录接口返回的 token，例如：Bearer eyJhbGciOi...
package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	settings "bluebell/setting"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	//1.加载配置
	if len(os.Args) < 2 {
		fmt.Printf("need config file.eg: bluebell config.yaml")
		return
	}
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success")
	//3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	//4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	if err := snowflake.Init(settings.Conf.StartTime, int64(settings.Conf.MachineID)); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	//初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}
	//5.注册路由
	r := routes.Setup(settings.Conf.Mode)
	//6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
