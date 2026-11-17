package routes

import (
	"bluebell/controller"
	_ "bluebell/docs"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-contrib/pprof"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}

	r := gin.New()
	//对全部限速
	//r.Use(logger.GinLogger(), logger.GinRecovery(true),middlewares.RateLimitMiddleware(2*time.Second,1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//在 Gin 中注册一个 Swagger 文档访问路由，让我们可以通过浏览器访问 API 接口文档，并在线测试接口。
	//加载 templates 文件夹下的 index.html 文件。
	r.LoadHTMLGlob("./templates/index.html")
	// /static：浏览器访问时使用的 URL 前缀。
	// ./static：服务器本地存放静态资源的文件夹。
	//浏览器访问：http://localhost:8080/static/css/style.css
	//Gin 就会去服务器本地寻找：./static/css/style.css然后将文件返回给浏览器。
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		//nil不向模板传递额外数据
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	//1.注册业务路由
	//从路由走到controller这一层，controller负责处理路由，参数校验，请求转发
	v1.POST("/signup", controller.SignUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)
	//根据时间或分数获取帖子列表，可查社区
	v1.GET("/posts2", controller.GetPostListHandler2)
	//应用JWT认证中间件,以下操作需要登录
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts/", controller.GetPostListHandler)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}
	pprof.Register(r) //注册pprof相关路由

	return r
}
