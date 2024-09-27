package middlewares

import (
	"blog/dao/redis"
	"blog/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveUserIp() func(c *gin.Context) {
	return func(c *gin.Context) {
		redis.SaveUserIp(c.ClientIP())
		c.Next() // 后续的处理请求的函数中 可以用过c.Get(CtxUserIDKey) 来获取当前请求的用户信息
	}
}

func IpLimit() func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := logic.IpLimit(c.ClientIP()); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"msg": fmt.Sprintf("%v", err),
			})
			c.Abort()
		}
		c.Next()
	}
}
