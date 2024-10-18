package controller

import (
	"blog/logic"
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

<<<<<<< HEAD:controller/data.go
func ResponseDataAboutIndexHandler(c *gin.Context) {
	var data = new(models.DataAboutIndex)
	err := logic.GetDataAboutIndex(data)
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, *data)
}

func ResponseDataAboutEssayListHandler(c *gin.Context) {
=======
func ResponseEssayListHandler(c *gin.Context) {
>>>>>>> dev:controller/essay_list.go
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

<<<<<<< HEAD:controller/data.go
	query.Classify = c.Query("classify")

	var essayListAndPage = new(models.DataAboutEssayListAndPage)
	essayListAndPage.EssayList = new([]models.DataAboutEssay)
	if err := logic.GetDataAboutClassifyEssayMsg(essayListAndPage, query); err != nil {
=======
	lID, _ := strconv.ParseInt(c.Query("label_id"), 10, 64)
	KID, _ := strconv.ParseInt(c.Query("kind_id"), 10, 64)

	query.LabelID = int(lID)
	query.KindID = int(KID)

	fmt.Println(query)

	var essayListAndPage = new(models.DataAboutEssayListAndPage)
	if err := logic.GetEssayList(essayListAndPage, query); err != nil {
>>>>>>> dev:controller/essay_list.go
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
