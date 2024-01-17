package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var r *gin.Engine

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
