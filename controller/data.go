package controller

import (
	"blog/dao/mysql"
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseDataAboutIndex(c *gin.Context) {
	//进入页面之后 得到各大分类种类以及相应的名称
	//1.查classifyKind和icon
	var classify = new([]models.DataAboutClassify)

	var err error
	err = mysql.GetDataAboutClassifyKind(classify)
	if err != nil {
		zap.L().Error("mysql.GetDataAboutClassifyKind(classify) failed", zap.Error(err))
		return
	}
	//2.查classifyDetails
	var classifyDetails = new([]models.DataAboutClassifyDetails)
	err = mysql.GetDataAboutClassifyDetails(classifyDetails)
	if err != nil {
		zap.L().Error("mysql.GetDataAboutClassifyDetails(classifyDetails) failed", zap.Error(err))
		return
	}

	//查essayDetail
	var essaysDetail = new([]models.DataAboutEssay)
	err = mysql.GetDataAboutClassifyEssayMsg(essaysDetail)
	// 3. 使用 map 优化循环 --> 整合数据
	//整合ClassifyName和ClassifyEssay
	var essaysDetailMap = make(map[string][]models.DataAboutEssay)

	for _, essay := range *essaysDetail {
		essaysDetailMap[essay.Kind] = append(essaysDetailMap[essay.Kind], models.DataAboutEssay{
			Name:         essay.Name,
			Kind:         essay.Kind,
			Router:       essay.Router,
			Introduction: essay.Introduction,
			ID:           essay.ID,
		})
	}
	//整合ClassifyKind和ClassifyName
	DataDetailMap := make(map[string][]models.DataAboutClassifyDetails)
	for _, detail := range *classifyDetails {
		if _, ok := essaysDetailMap[detail.Name]; ok { //得到了该分类的所有文章
			data := models.DataAboutClassifyDetails{
				Kind:   detail.Kind,
				Name:   detail.Name,
				Router: detail.Router,
				Essay:  essaysDetailMap[detail.Name],
				ID:     detail.ID,
			}

			DataDetailMap[detail.Kind] = append(DataDetailMap[detail.Kind], data)
		}
	}
	// 4. 初始化 DataAboutIndex 切片
	DataAboutIndexMenu := make([]models.DataAboutIndexMenu, len(*classify))

	// 5. 将 classify 和 classify 相对应写成【map，map】格式
	for i := 0; i < len(*classify); i++ {
		DataAboutIndexMenu[i].ClassifyKind = (*classify)[i].ClassifyKind
		DataAboutIndexMenu[i].Icon = (*classify)[i].Icon
		DataAboutIndexMenu[i].ClassifyDetails = DataDetailMap[(*classify)[i].ClassifyKind]
	}

	//6.得到该用户的信息
	var userInfo = new(models.UserInfo)

	var DataAboutIndex = models.DataAboutIndex{
		DataAboutIndexMenu: DataAboutIndexMenu,
		UserInfo:           *userInfo,
	}
	ResponseSuccess(c, DataAboutIndex)
}

func ResponseDataAboutEssay(c *gin.Context) {
	//1.参数处理
	query := c.Query("id")
	id, _ := strconv.Atoi(query)
	//2.业务处理
	var essay = new(models.EssayContent)
	err := logic.GetEssayData(essay, id)
	if err != nil {
		zap.L().Error("logic.GetEssayData(essay, id) failed", zap.Error(err))
		return
	}
	//3.返回响应
	ResponseSuccess(c, essay)
}
