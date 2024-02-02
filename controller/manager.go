package controller

import (
	"blog/dao/mysql"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	addClassifySuccess = "添加分类成功"
	addEssay           = "添加文章成功"
)

func AddClassifyHandler(c *gin.Context) {
	//1.参数处理
	var classify = new(models.ClassifyParams)
	if err := c.ShouldBindJSON(classify); err != nil {
		zap.L().Error("c.ShouldBindJSON(classify) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	err := mysql.AddClassify(classify)
	if err != nil {
		zap.L().Error("mysql.AddClassify(classify) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, addClassifySuccess)
}

func AddEssayHandler(c *gin.Context) {
	//1.参数处理
	var essay = new(models.EssayParams)
	if err := c.ShouldBindJSON(essay); err != nil {
		zap.L().Error("c.ShouldBindJSON(essay) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//2.业务处理
	if err := mysql.CreateEssay(essay); err != nil {
		zap.L().Error("mysql.CreateEssay(essay) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, addEssay)
}
