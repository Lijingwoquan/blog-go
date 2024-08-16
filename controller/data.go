package controller

import (
	"blog/cache"
	"blog/dao/mysql"
	"blog/logic"
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseDataAboutIndexHandler(c *gin.Context) {
	p := c.Query("page")
	fmt.Println(p)
	//得到各大分类种类以及相应的名称
	ResponseSuccess(c, cache.GlobalDataAboutIndex)
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

func ResponseDataAboutIndexAsideHandler(c *gin.Context) {
	if cache.Error != nil {
		ResponseError(c, CodeServeBusy)
	}
	ResponseSuccess(c, cache.GlobalDataAboutIndex)
}

// 1.查kind和icon
func getKindAndIcon(k *[]models.DataAboutKind) error {
	return mysql.GetDataAboutKind(k)
}

// 2.查classifyDetails
func getClassifyAndDetails(c *[]models.DataAboutClassify) error {
	//得到了所有的分类
	return mysql.GetDataAboutClassifyDetails(c)
}

// 3.整合kind数据和classify数据
func sortIndexData(k *[]models.DataAboutKind, c *[]models.DataAboutClassify) *models.DataAboutIndex {
	var indexData = models.DataAboutIndex{}
	var indexDataMenu = make([]models.DataAboutIndexMenu, len(*k))

	var kindMap = make(map[string][]models.DataAboutKind)
	var kindAndClassifyMap = make(map[string][]models.DataAboutClassify)

	// 自下而上 --> 先遍历c 然后
	for _, classify := range *c {
		kindAndClassifyMap[classify.Kind] = append(kindAndClassifyMap[classify.Name], classify)
	}

	for i, kind := range *k {
		indexDataMenu[i].DataAboutKind = kind
		kindMap[kind.ClassifyKind] = append(kindMap[kind.ClassifyKind], kind)
	}

	for i, kind := range *k {
		indexDataMenu[i].Classify = kindAndClassifyMap[kind.ClassifyKind]
	}

	indexData.DataAboutIndexMenu = indexDataMenu
	return &indexData
}
