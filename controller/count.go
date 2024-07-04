package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUserIpCountHandel(c *gin.Context) {
	setKind := new(models.SetKind)
	if err := logic.GetUserIpCount(setKind); err != nil {
		zap.L().Error("redis.GetUserIpCount(setKind) failed,err", zap.Error(err))
		return
	}
	ResponseSuccess(c, setKind)
}
