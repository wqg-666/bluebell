package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	//fillInterval	令牌补充的时间间隔
	//cap	令牌桶的最大容量
	bucket := ratelimit.NewBucket(fillInterval, cap) //创建令牌桶
	//创建之后，bucket 就代表这个令牌桶。
	//注意：这里的 bucket 在中间件外层创建，而不是在里面创建。
	//这意味着什么？
	//所有经过这个中间件的请求，共享同一个令牌桶。
	//如果每次请求都重新创建令牌桶，那么每个请求都有自己的桶，限流就失去意义了。
	return func(c *gin.Context) {
		// 如果取不到令牌就返回响应
		//TakeAvailable(1) 的意思是： 尝试从令牌桶中取出 1 个令牌。它的返回值表示实际成功获取的令牌数量
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		//取到令牌就放行
		c.Next()
	}
}
