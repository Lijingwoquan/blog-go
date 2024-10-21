package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	createLabelSuccess    = "添加分类成功"
	updateClassifySuccess = "修改分类成功"
)

func CreateLabelHandler(c *gin.Context) {
	//1.参数处理
	var label = new(models.LabelParams)
	if err := c.ShouldBindJSON(label); err != nil {
		zap.L().Error("c.ShouldBindJSON(label) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.CreateLabel(label); err != nil {
		zap.L().Error("mysql.CreateLabel(classify) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, createLabelSuccess)
}

func DeleteLabelHandler(c *gin.Context) {

}

func UpdateLabelHandler(c *gin.Context) {
	//1.参数处理
	var label = new(models.LabelParams)
	if err := c.ShouldBindJSON(label); err != nil {
		zap.L().Error("c.ShouldBindJSON(label) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.UpdateLabel(label); err != nil {
		zap.L().Error("logic.UpdateClassify(classify) failed err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateClassifySuccess)
}
