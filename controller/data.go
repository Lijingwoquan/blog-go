package controller

import (
	"blog/cache"
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseDataAboutIndexAsideHandler(c *gin.Context) {
	if cache.Error != nil {
		ResponseError(c, CodeServeBusy)
	}
	ResponseSuccess(c, cache.GlobalDataAboutIndex)
}

func ResponseDataAboutEssayListHandler(c *gin.Context) {
	query := models.EssayQuery{}
	page64, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	if page64 == 0 {
		page64 = 1
	}
	query.Page = int(page64)

	pageSize64, _ := strconv.ParseInt(c.Query("pageSize"), 10, 64)
	if pageSize64 == 0 {
		pageSize64 = 5
	}
	query.PageSize = int(pageSize64)

	query.Classify = c.Query("classify")

	var essayListAndPage = new(models.DataAboutEssayListAndPage)
	essayListAndPage.EssayList = new([]models.DataAboutEssay)
	if err := logic.GetDataAboutClassifyEssayMsg(essayListAndPage, query); err != nil {
		zap.L().Error("logic.GetDataAboutClassifyEssayMsg(essayList) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//得到各大分类种类以及相应的名称
	ResponseSuccess(c, essayListAndPage)
}

func ResponseDataAboutEssayHandler(c *gin.Context) {
	//1.参数处理
	query := c.Query("id")
	id, err := strconv.Atoi(query)
	if err != nil {
		zap.L().Error("strconv.Atoi(query) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//2.业务处理
	var essay = new(models.EssayData)
	if err = logic.GetEssayData(essay, id); err != nil {
		zap.L().Error("logic.GetEssayData(essay, id) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, essay)
}
