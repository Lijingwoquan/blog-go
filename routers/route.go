package routers

import (
	"blog/controller"
	"blog/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	//这里要设置端口的 前端是:80不用显示调用
	config.AllowOrigins = []string{"https://liuzihao.online", "https://www.liuzihao.online", "http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	r.Use(cors.New(config))
	r.Static("api/img", "/app/statics/img")
	r.Static("api/file", "/app/statics/file")

	v1 := r.Group("/api/base")
	{
		// 使用中间件的路由
		v1.GET("/index", controller.ResponseDataAboutIndexHandler)
		v1.GET("/essay", controller.ResponseDataAboutEssayHandler)
	}

	v2 := r.Group("/api/user")
	{
		v2.POST("/login", controller.LoginHandler)
		//v2.POST("/signup", controller.SignupHandler)
		//v2.POST("/logout", controller.LogoutHandler)
		//v2.POST("/updateUserMsg", middlewares.JWTAuthMiddleware(), controller.UpdateUserMsgHandler)
	}

	v3 := r.Group("/api/manager")
	v3.Use(middlewares.JWTAuthMiddleware(), middlewares.UpdateDataMiddleware())
	{
		v3.POST("/addClassify", controller.AddClassifyHandler)
		v3.POST("/addEssay", controller.AddEssayHandler)
		v3.PUT("/updateKind", controller.UpdateKindHandler)
		v3.PUT("/updateClassify", controller.UpdateClassifyHandler)
		v3.PUT("/updateEssayMsg", controller.UpdateEssayMSgHandler)
		v3.PUT("/updateEssayContent", controller.UpdateEssayContentHandler)
		v3.DELETE("/deleteEssay", controller.DeleteEssayHandler)
		v3.POST("/uploadImg", controller.UploadImgHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSONP(404, gin.H{
			"msg": "404",
		})
	})
	return r
}
