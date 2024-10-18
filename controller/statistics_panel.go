package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ResponseDataAboutManagerPanel(c *gin.Context) {
	panelList := new(models.Panel)

	if err := logic.GetUserIpCount(&panelList.IpSet); err != nil {
		zap.L().Error("logic.GetUserIpCount(&panelList.IpSet) failed,err", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	// 逻辑处理
	if err := logic.GetSearchKeywordRank(&panelList.RankZset); err != nil {
		zap.L().Error(" logic.GetSearchKeywordRank(&panelList.RankZset) failed,err:", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//	返回响应
	ResponseSuccess(c, panelList)
}
