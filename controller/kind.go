package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	createKindSuccess = "创建kind成功"
	deleteKindSuccess = "删除kind成功"
	updateKindSuccess = "修改kind成功"
)

func CreateKindHandler(c *gin.Context) {
	k := new(models.KindParams)
	// 1.参数绑定
	if err := c.ShouldBindJSON(k); err != nil {
		zap.L().Error("c.ShouldBindJSON(k) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	// 2.逻辑处理
	if err := logic.CreateKind(k); err != nil {
		zap.L().Error("logic.CreateKind(k) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, createKindSuccess)
}

func DeleteKindHandler(c *gin.Context) {

}

func UpdateKindHandler(c *gin.Context) {
	//1.参数检验
	var k = new(models.KindUpdateParams)
	if err := c.ShouldBindJSON(k); err != nil {
		zap.L().Error("c.ShouldBindJSON(k) failed,err:", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}
	//2.逻辑处理
	if err := logic.UpdateKind(k); err != nil {
		zap.L().Error("logic.UpdateKind(k) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, updateKindSuccess)
}
