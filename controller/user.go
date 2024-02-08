package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

const (
	signupSuccess        = "注册成功"
	userIDInvalid        = "无法获取该用户id"
	updateUserMsgSuccess = "修改个人信息成功"
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

// LogoutHandler 退出登录
func LogoutHandler(c *gin.Context) {
	//1.参数验证 --> 得到相应的token
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	//得到token
	token := parts[1]

	//2.业务处理 --> 将该token储存在数据库中
	err := logic.Logout(token)
	if err != nil {
		zap.L().Error("logic.Logout(token) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, CodeSuccess)
}

// UpdateUserMsgHandler 修改用户信息
func UpdateUserMsgHandler(c *gin.Context) {
	//1.参数校验
	var user = new(models.UserParams)
	err := c.ShouldBindJSON(user)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON(user) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	id := getUserId(c)

	//2.业务处理
	err = logic.UpdateUserMsg(user, id)
	if err != nil {
		zap.L().Error("logic.UpdateUserMsg(fmt.Sprintf(\"%s\",id)) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, updateUserMsgSuccess)
}
