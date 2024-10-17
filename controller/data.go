package controller

import (
	"blog/logic"
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseIndexDataHandler(c *gin.Context) {
	var data = new(models.DataAboutIndex)
	var err error
	data, err = logic.GetIndexData()
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, *data)
}

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

	fmt.Println(query)

	var essayListAndPage = new(models.DataAboutEssayListAndPage)
	if err := logic.GetEssayList(essayListAndPage, query); err != nil {
		zap.L().Error("logic.GetDataAboutClassifyEssayMsg(essayList) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//得到各大分类种类以及相应的名称
	ResponseSuccess(c, essayListAndPage)
}

func ResponseDataAboutEssayHandler(c *gin.Context) {
	//1.参数处理
	queryID := c.Query("id")
	id, err := strconv.Atoi(queryID)
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
