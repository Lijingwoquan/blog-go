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
		classifyNameIncludeEssayMap[include.Kind] = append(classifyNameIncludeEssayMap[include.Kind], models.ClassifyIncludeEssay{
			Name:   include.Name,
			Kind:   include.Kind,
			Router: include.Router,
		})
	}
	//整合ClassifyKind和ClassifyName
	DataDetailMap := make(map[string][]models.DataAboutClassifyDetails)
	for _, detail := range *classifyDetails {
		if _, ok := classifyNameIncludeEssayMap[detail.Name]; ok { //得到了该分类的所有文章
			data := models.DataAboutClassifyDetails{
				Kind:   detail.Kind,
				Name:   detail.Name,
				Router: detail.Router,
				Essay:  classifyNameIncludeEssayMap[detail.Name],
			}

			DataDetailMap[detail.Kind] = append(DataDetailMap[detail.Kind], data)
		}
	}
	// 4. 初始化 DataAboutIndex 切片
	DataAboutIndexMenu := make([]models.DataAboutIndexMenu, len(*classifyKind))

	// 5. 将 classify 和 classify 相对应写成【map，map】格式
	for i := 0; i < len(*classifyKind); i++ {
		DataAboutIndexMenu[i].ClassifyKind = (*classifyKind)[i]
		DataAboutIndexMenu[i].ClassifyDetails = DataDetailMap[(*classifyKind)[i]]
	}

	//6.得到该用户的信息
	var userInfo = new(models.UserInfo)

	var DataAboutIndex = models.DataAboutIndex{
		DataAboutIndexMenu: DataAboutIndexMenu,
		UserInfo:           *userInfo,
	}
	ResponseSuccess(c, DataAboutIndex)
}
