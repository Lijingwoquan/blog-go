package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	updateKindSuccess = "修改大纲成功"
)

func CreateKindHandler(c *gin.Context) {
	//if err := c.ShouldBindJSON(); err != nil {
	//
	//}
}

func DeleteKindHandler(c *gin.Context) {

}

func UpdateKindHandler(c *gin.Context) {
	//1.参数检验
	var k = new(models.KindParams)
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
