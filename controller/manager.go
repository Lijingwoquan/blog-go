package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	addClassifySuccess = "添加分类成功"
	addEssaySuccess    = "添加文章成功"
	updateEssaySuccess = "修改文章成功"
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
	err := logic.AddClassify(classify)
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
	if err := logic.CreateEssay(essay); err != nil {
		zap.L().Error("mysql.CreateEssay(essay) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, addEssaySuccess)
}

func UpdateEssayHandler(c *gin.Context) {
	//1.获取参数
	var data = new(models.UpdateEssay)
	if err := c.ShouldBindJSON(data); err != nil {
		zap.L().Error("c.ShouldBindJSON(data) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.UpdateEssay(data); err != nil {
		zap.L().Error("mysql.UpdateEssay(data) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateEssaySuccess)
}
