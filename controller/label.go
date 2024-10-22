package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const (
	createLabelSuccess = "创建label成功"
	updateLabelSuccess = "修改label成功"
	deleteLabelSuccess = "删除label成功"
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
	//1.参数处理
	ids := c.Query("id")
	id64, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		zap.L().Error("parse id failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	id := int(id64)
	//2.业务处理
	if err := logic.DeleteLabel(id); err != nil {
		zap.L().Error("logic.DeleteLabel(id) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, deleteLabelSuccess)
}

func UpdateLabelHandler(c *gin.Context) {
	//1.参数处理
	var label = new(models.LabelUpdateParams)
	if err := c.ShouldBindJSON(label); err != nil {
		zap.L().Error("c.ShouldBindJSON(label) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.业务处理
	if err := logic.UpdateLabel(label); err != nil {
		zap.L().Error("logic.UpdateLabel(label) failed err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateLabelSuccess)
}
