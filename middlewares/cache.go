package middlewares

import (
	"blog/cache"
	"github.com/gin-gonic/gin"
)

func UpdateDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 在请求被处理之前，不做任何事情

		// 调用下一个中间件或处理函数
		c.Next()

		// 在请求被处理之后，更新数据
		cache.UpdateDataAboutIndex()
	}
}
