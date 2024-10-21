package help

import (
	"blog/dao/mysql"
	"blog/models"
)

<<<<<<< HEAD
// ResponseDataAboutIndex 返回首页数据
func ResponseDataAboutIndex(DataAboutIndex *models.DataAboutIndex) (err error) {
	var kindList = new([]models.DataAboutKind)
	if err = getKindAndIcon(kindList); err != nil {
		return err
	}
	var classifyList = new([]models.DataAboutClassify)
	if err = getClassifyAndDetails(classifyList); err != nil {
=======
// ResponseIndexData 返回首页数据
func ResponseIndexData(DataAboutIndex *models.IndexData) (err error) {
	var kindList = new([]models.KindData)
	if err = mysql.GetKindList(kindList); err != nil {
		return err
	}
	var labelList = new([]models.LabelData)
	if err = mysql.GetLabelList(labelList); err != nil {
		return err
	}

	var essayList = new([]models.EssayData)
	if err = mysql.GetRecommendEssayList(essayList); err != nil {
>>>>>>> dev
		return err
	}

	//整合数据
	sortKindAndClassify(DataAboutIndex, kindList, classifyList)

	return nil
}
<<<<<<< HEAD:help/data.go

func ResponseDataAboutEssayList(essayList *[]models.DataAboutEssay) (err error) {
	if err = getAllEssay(essayList); err != nil {
		return err
	}
	// 整合classify
	var classifyList = new([]models.DataAboutClassify)
	if err = mysql.GetAllDataAboutClassify(classifyList); err != nil {
		return err
	}
	var classifyMap = make(map[string]string)
	for _, classify := range *classifyList {
		classifyMap[classify.Name] = classify.Router
	}
	// 再遍历essayList为它们加上kindRouter
	for i, essay := range *essayList {
		(*essayList)[i].KindRouter = classifyMap[essay.Kind]
	}
	return redis.GetEssayKeywordsForIndex(essayList)
}

// 1.查kind和icon
func getKindAndIcon(k *[]models.DataAboutKind) error {
	return mysql.GetDataAboutKind(k)
}

// 2.查classifyDetails
func getClassifyAndDetails(c *[]models.DataAboutClassify) error {
	//得到了所有的分类
	return mysql.GetAllDataAboutClassify(c)
}

// 3.查询allEssay
func getAllEssay(data *[]models.DataAboutEssay) error {
	return mysql.GetAllEssay(data)
}

// 4.整合数据
func sortKindAndClassify(DataAboutIndex *models.DataAboutIndex, k *[]models.DataAboutKind, c *[]models.DataAboutClassify) {
	var indexDataMenu = make([]models.DataAboutIndexMenu, len(*k))

	var kindAndClassifyMap = make(map[string][]models.DataAboutClassify)
	for _, classify := range *c {
		kindAndClassifyMap[classify.Kind] = append(kindAndClassifyMap[classify.Kind], classify)
	}

	for i, kind := range *k {
		indexDataMenu[i].Kind = kind
		indexDataMenu[i].ClassifyList = kindAndClassifyMap[kind.Name]
	}

	(*DataAboutIndex).Menu = indexDataMenu
}
=======
>>>>>>> dev:help/index_data.go
