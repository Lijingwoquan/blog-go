package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseEssayListHandler(c *gin.Context) {
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

	lID, _ := strconv.ParseInt(c.Query("label_id"), 10, 64)
	KID, _ := strconv.ParseInt(c.Query("kind_id"), 10, 64)

	query.LabelID = int(lID)
	query.KindID = int(KID)

	var essayListAndPage = new(models.DataAboutEssayListAndPage)
	if err := logic.GetEssayList(essayListAndPage, query); err != nil {
		zap.L().Error("logic.GetDataAboutClassifyEssayMsg(essayList) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//得到各大分类种类以及相应的名称
	ResponseSuccess(c, essayListAndPage)
}

func ResponseDataAboutSearchKeyword(c *gin.Context) {
	//1.参数检验
	searchParam := new(models.SearchParam)
	if err := c.ShouldBindJSON(searchParam); err != nil {
		zap.L().Error("c.ShouldBindJSON(keyword) failed", zap.Error(err))
		ResponseError(c, CodeParamInvalid)
		return
	}

	var essayList = new([]models.DataAboutEssay)
	//2.逻辑处理
	if err := logic.GetDataByKeyword(essayList, searchParam); err != nil {
		zap.L().Error("logic.IncreaseSearchKeyword(keyword) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, essayList)
}
