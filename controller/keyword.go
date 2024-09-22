package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ResponseDataAboutSearchKeyword(c *gin.Context) {
	//1.参数检验
	searchParam := new(models.SearchParam)
	if err := c.ShouldBindJSON(searchParam); err != nil {
		zap.L().Error("c.ShouldBindJSON(keyword) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}

	//2.逻辑处理
	if err := logic.GetDataByKeyword(searchParam); err != nil {
		zap.L().Error("logic.IncreaseSearchKeyword(keyword) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//3.返回响应
	var essayList = new([]models.DataAboutEssay)
	logic.GetEssayListByKeyword(essayList, searchParam)
	ResponseSuccess(c, essayList)
}
