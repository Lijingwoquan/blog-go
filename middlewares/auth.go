package middlewares

import (
	"blog/controller"
	"blog/dao/mysql"
	"blog/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxxx.xxx.xxx  / X-TOKEN: xxx.xxx.xx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			//controller.ResponseError(c, controller.CodeNeedLogin)
			c.JSON(http.StatusOK, gin.H{
				"msg": "需要登录",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"msg": "需要登录",
			})
			c.Abort()
			return
		}
		if err := JWTInvalidToken(parts[1]); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "需要登录",
			})
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "需要登录",
			})
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理请求的函数中 可以用过c.Get(CtxUserIDKey) 来获取当前请求的用户信息
	}
}

func JWTInvalidToken(token string) (err error) {
	//1.查看token是否在数据库中
	err = mysql.CheckTokenIfInvalid(token)
	return err
}
