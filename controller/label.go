package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	addClassifySuccess    = "添加分类成功"
	updateClassifySuccess = "修改分类成功"
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
	if err := logic.AddClassify(classify); err != nil {
		zap.L().Error("mysql.AddClassify(classify) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, addClassifySuccess)
}

func UpdateClassifyHandler(c *gin.Context) {
	//1.参数处理
	var classify = new(models.UpdateClassifyParams)
	if err := c.ShouldBindJSON(classify); err != nil {
		zap.L().Error("c.ShouldBindJSON(classify) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.UpdateClassify(classify); err != nil {
		zap.L().Error("logic.UpdateClassify(classify) failed err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateClassifySuccess)
}
