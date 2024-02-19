package controller

import (
	"blog/dao/mysql"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseDataAboutIndexHandler(c *gin.Context) {
	//进入页面之后 得到各大分类种类以及相应的名称
	//1.查kind和icon
	var kindDetails = new([]models.DataAboutKind)

	var err error
	err = mysql.GetDataAboutKind(kindDetails)
	if err != nil {
		zap.L().Error("mysql.GetDataAboutKind(kindDetails) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//2.查classifyDetails
	var classifyDetails = new([]models.DataAboutClassify)
	err = mysql.GetDataAboutClassifyDetails(classifyDetails) //得到了所有的分类
	if err != nil {
		zap.L().Error("mysql.GetDataAboutClassifyDetails(classifyDetails) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//查essayDetail
	var essaysDetail = new([]models.DataAboutEssay)
	err = mysql.GetDataAboutClassifyEssayMsg(essaysDetail)
	if err != nil {
		zap.L().Error("mysql.GetDataAboutClassifyEssayMsg(essaysDetail) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	// 3. 使用 map 优化循环 --> 整合数据
	//整合classify和essay
	var essaysDetailMap = make(map[string][]models.DataAboutEssay)
	//初始化map
	for _, classify := range *classifyDetails {
		essaysDetailMap[classify.Name] = []models.DataAboutEssay{} // 初始化为空的切片
	}
	for _, essay := range *essaysDetail {
		essaysDetailMap[essay.Kind] = append(essaysDetailMap[essay.Kind], models.DataAboutEssay{
			Name:         essay.Name,
			Kind:         essay.Kind,
			Router:       essay.Router,
			Introduction: essay.Introduction,
			ID:           essay.ID,
			CreatedTime:  essay.CreatedTime,
		})
	}

	//为分类好的essay加上page
	for _, valueArr := range essaysDetailMap {
		var page = 1
		var count = 0
		for i, v := range valueArr {
			count++
			if count > 10 {
				page++
				count = 0
			}
			v.Page = page
			valueArr[i].Page = page
		}
	}

	//整合Kind和Classify
	KindAndClassifyMap := make(map[string][]models.DataAboutClassify)
	for _, detail := range *classifyDetails {
		if _, ok := essaysDetailMap[detail.Name]; ok { //得到了该分类的所有文章
			data := models.DataAboutClassify{
				Kind:   detail.Kind,
				Name:   detail.Name,
				Router: detail.Router,
				Essay:  essaysDetailMap[detail.Name],
				ID:     detail.ID,
			}
			KindAndClassifyMap[detail.Kind] = append(KindAndClassifyMap[detail.Kind], data)
		}

	}
	// 4. 初始化 DataAboutIndex 切片
	DataAboutIndexMenu := make([]models.DataAboutIndexMenu, len(*kindDetails))

	// 5. 将 kind 和 classify 相对应写成【map，map】格式
	for i := 0; i < len(*kindDetails); i++ {
		DataAboutIndexMenu[i].ClassifyKind = (*kindDetails)[i].ClassifyKind
		DataAboutIndexMenu[i].Icon = (*kindDetails)[i].Icon
		DataAboutIndexMenu[i].Id = (*kindDetails)[i].Id
		DataAboutIndexMenu[i].Classify = KindAndClassifyMap[(*kindDetails)[i].ClassifyKind]
	}

	//7.整合数据
	var DataAboutIndex = models.DataAboutIndex{
		DataAboutIndexMenu: DataAboutIndexMenu,
	}
	ResponseSuccess(c, DataAboutIndex)
}

func ResponseDataAboutEssayHandler(c *gin.Context) {
	//1.参数处理
	query := c.Query("id")
	id, _ := strconv.Atoi(query)
	//2.业务处理
	var essay = new(models.EssayData)
	err := mysql.GetEssayData(essay, id)
	if err != nil {
		zap.L().Error("logic.GetEssayData(essay, id) failed", zap.Error(err))
		return
	}
	//3.返回响应
	ResponseSuccess(c, essay)
}
