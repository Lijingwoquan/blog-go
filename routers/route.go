package routers

import (
	"blog/controller"
	"blog/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

var r *gin.Engine

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.POST("/signup", controller.SignupHandler)
	r.POST("/login", controller.LoginHandler)

	r.Use(middlewares.JWTAuthMiddleware())
	{
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		})
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
