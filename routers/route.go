package routers

import (
	"blog/controller"
	"blog/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()
	r := gin.Default()
	//r.Use(cors.Default()) --> 这里没有Authorization！！！妈的被坑惨了

	// 创建新的CORS中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	v1 := r.Group("/api/base")
	{
		// 使用中间件的路由
		v1.GET("/index", middlewares.JWTAuthMiddleware(), controller.ResponseDataAboutIndexHandler)
		v1.GET("/essay", middlewares.JWTAuthMiddleware(), controller.ResponseDataAboutEssayHandler)
	}

	v2 := r.Group("/api/user")
	{
		v2.POST("/signup", controller.SignupHandler)
		v2.POST("/login", controller.LoginHandler)
		v2.POST("/logout", middlewares.JWTAuthMiddleware(), controller.LogoutHandler)
		v2.POST("/updateUserMsg", middlewares.JWTAuthMiddleware(), controller.UpdateUserMsgHandler)
	}

	v3 := r.Group("/api/manager")
	{
		v3.POST("/addClassify", controller.AddClassifyHandler)
		v3.POST("/addEssay", controller.AddEssayHandler)
		v3.PUT("/updateEssay", controller.UpdateEssayHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
