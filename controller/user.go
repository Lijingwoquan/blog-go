package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	signupSuccess = "注册成功"
	loginSuccess  = "登录成功"
)

// SignupHandler 注册
func SignupHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var u = new(models.UserParams)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("c.ShouldBindJSON(&u) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.Signup(u); err != nil {
		ResponseError(c, CodeUserExist)
		return
	}
	//3.返回响应
	ResponseSuccess(c, signupSuccess)
}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	//1.获取参数并检验
	var u = new(models.User)
	if err := c.ShouldBindJSON(u); err != nil {
		zap.L().Error("c.ShouldBindJSON(u) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.Login(u); err != nil {
		zap.L().Error("logic.Login() failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, u.Token)
}
