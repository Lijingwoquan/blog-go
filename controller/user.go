package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SignHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var u = new(models.UserParams)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("c.ShouldBindJSON(&u) failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
		return
	}
	//2.业务处理
	if err := logic.Signup(u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err,
		})
		return
	}
	//3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "注册成功",
	})
}
func LoginHandler(c *gin.Context) {
	//1.获取参数并检验
	var u = new(models.User)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("c.ShouldBindJSON(u) failed", zap.Error(err))
		return
	}
	//2.业务处理
	if err := logic.Login(u); err != nil {
		zap.L().Error("logic.Login() failed", zap.Error(err))
		return
	}

	//3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": u,
	})
}
