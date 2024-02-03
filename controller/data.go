package controller

import (
	"blog/dao/mysql"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ResponseDataAboutIndex(c *gin.Context) {
	//进入页面之后 得到各大分类种类以及相应的名称
	//1.查classifyKind
	var classifyKind = new([]string)
	var err error
	err = mysql.GetDataAboutClassifyKind(classifyKind)
	if err != nil {
		zap.L().Error("mysql.GetDataAboutClassifyKind() failed", zap.Error(err))
		return
	}
	//2.查DataAboutClassifyDetails
	var classifyDetails = new([]models.DataAboutClassifyDetails)
	err = mysql.GetDataAboutClassifyDetails(classifyDetails)
	if err != nil {
		zap.L().Error("mysql.mysql.GetDataAboutClassifyKind() failed", zap.Error(err))
		return
	}

	//查ClassifyEssayName
	var classifyNameIncludeEssay = new([]models.ClassifyIncludeEssay)
	err = mysql.GetDataAboutClassifyEssayName(classifyNameIncludeEssay)

	// 3. 使用 map 优化循环 --> 整合数据
	//整合ClassifyName和ClassifyEssay
	var classifyNameIncludeEssayMap = make(map[string][]models.ClassifyIncludeEssay)

	for _, include := range *classifyNameIncludeEssay {
		classifyNameIncludeEssayMap[include.EssayKind] = append(classifyNameIncludeEssayMap[include.EssayKind], models.ClassifyIncludeEssay{
			EssayName:  include.EssayName,
			EssayKind:  include.EssayKind,
			EssayRoute: include.EssayRoute,
		})
	}
	//整合ClassifyKind和ClassifyName
	DataDetailMap := make(map[string][]models.DataAboutClassifyDetails)
	for _, detail := range *classifyDetails {
		if _, ok := classifyNameIncludeEssayMap[detail.ClassifyName]; ok { //得到了该分类的所有文章
			data := models.DataAboutClassifyDetails{
				ClassifyKindName: detail.ClassifyKindName,
				ClassifyName:     detail.ClassifyName,
				ClassifyRoute:    detail.ClassifyRoute,
				ClassifyEssay:    classifyNameIncludeEssayMap[detail.ClassifyName],
			}

			DataDetailMap[detail.ClassifyKindName] = append(DataDetailMap[detail.ClassifyKindName], data)
		}
	}
	// 4. 初始化 DataAboutIndex 切片
	DataAboutIndex := make([]models.DataAboutIndex, len(*classifyKind))

	// 5. 将 classify 和 classify 相对应写成【map，map】格式
	for i := 0; i < len(*classifyKind); i++ {
		DataAboutIndex[i].ClassifyKindName = (*classifyKind)[i]
		DataAboutIndex[i].ClassifyDetails = DataDetailMap[(*classifyKind)[i]]
	}
	ResponseSuccess(c, DataAboutIndex)
}
