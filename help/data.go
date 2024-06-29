package help

import (
	"blog/dao/mysql"
	"blog/models"
)

const (
	maxEssaySumForPage = 5
)

// ResponseDataAboutIndex 返回首页数据
func ResponseDataAboutIndex(DataAboutIndex *models.DataAboutIndex) (err error) {
	//进入页面之后 得到各大分类种类以及相应的名称
	//kind-->classify-->essay

	//1.查kind和icon
	var kindDetails = new([]models.DataAboutKind)
	if err = getKindAndIcon(kindDetails); err != nil {
		return
	}
	//2.查classifyDetails
	var classifyDetails = new([]models.DataAboutClassify)
	if err = getClassifyAndDetails(classifyDetails); err != nil {
		return
	}
	//3.查essayDetail
	var essaysDetail = new([]models.DataAboutEssay)
	if err = getEssayDetail(essaysDetail); err != nil {
		return
	}

	//4. 整合classify和essay 使用 map 优化循环
	var classifyAndEssayMap = make(map[string][]models.DataAboutEssay, 10)
	sortClassifyAndEssay(classifyAndEssayMap, classifyDetails, essaysDetail)

	//5.整合Kind和Classify
	kindAndClassifyMap := make(map[string][]models.DataAboutClassify, 10)
	sortKindAndClassify(kindAndClassifyMap, classifyAndEssayMap, classifyDetails)

	//6. 将 kind 和 classify 相对应写成【map，map】格式
	//初始化 DataAboutIndex 切片
	DataAboutIndexMenu := make([]models.DataAboutIndexMenu, len(*kindDetails))
	for i := 0; i < len(*kindDetails); i++ {
		DataAboutIndexMenu[i].ClassifyKind = (*kindDetails)[i].ClassifyKind
		DataAboutIndexMenu[i].Icon = (*kindDetails)[i].Icon
		DataAboutIndexMenu[i].Id = (*kindDetails)[i].Id
		DataAboutIndexMenu[i].Classify = kindAndClassifyMap[(*kindDetails)[i].ClassifyKind]
	}

	//7.返回数据处理
	*DataAboutIndex = models.DataAboutIndex{
		DataAboutIndexMenu: DataAboutIndexMenu,
	}
	return nil
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

// 3.查essayDetail
func getEssayDetail(c *[]models.DataAboutEssay) error {
	return mysql.GetDataAboutClassifyEssayMsg(c)
}

// 4.整合Kind和Classify
func sortClassifyAndEssay(classifyAndEssayMap map[string][]models.DataAboutEssay, classifyDetails *[]models.DataAboutClassify, essaysDetail *[]models.DataAboutEssay) {
	//初始化map 给每个分类一个空的切片 这个切片用来存放该分类下的文章
	for _, classify := range *classifyDetails {
		classifyAndEssayMap[classify.Name] = []models.DataAboutEssay{} // 初始化为空的切片
	}
	//将文章按照分类放入map中
	for _, essay := range *essaysDetail {
		classifyAndEssayMap[essay.Kind] = append(classifyAndEssayMap[essay.Kind], models.DataAboutEssay{
			Name:         essay.Name,
			Kind:         essay.Kind,
			Router:       essay.Router,
			Introduction: essay.Introduction,
			ID:           essay.ID,
			CreatedTime:  essay.CreatedTime,
		})
	}
	//为分类好的essay加上page
	for _, essayArr := range classifyAndEssayMap {
		var page = 1
		var count = 0
		for i, v := range essayArr {
			count++
			if count > maxEssaySumForPage {
				page++
				count = 0
			}
			v.Page = page
			essayArr[i].Page = page
		}
	}
}

// 5.整合Kind和Classify
func sortKindAndClassify(kindAndClassifyMap map[string][]models.DataAboutClassify, classifyAndEssayMap map[string][]models.DataAboutEssay, classifyDetails *[]models.DataAboutClassify) {
	for _, detail := range *classifyDetails {
		if _, ok := classifyAndEssayMap[detail.Name]; ok { //得到了该分类的所有文章
			data := models.DataAboutClassify{
				Kind:   detail.Kind,
				Name:   detail.Name,
				Router: detail.Router,
				Essay:  classifyAndEssayMap[detail.Name],
				ID:     detail.ID,
			}
			kindAndClassifyMap[detail.Kind] = append(kindAndClassifyMap[detail.Kind], data)
		}
	}
}
