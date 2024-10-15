package help

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

// ResponseDataAboutIndex 返回首页数据
func ResponseDataAboutIndex(DataAboutIndex *models.DataAboutIndex) (err error) {
	var kindList = new([]models.DataAboutKind)
	if err = getKindAndIcon(kindList); err != nil {
		return err
	}
	var classifyList = new([]models.DataAboutLabel)
	if err = getClassifyAndDetails(classifyList); err != nil {
		return err
	}

	//整合数据
	DataAboutIndex.KindList = *kindList
	DataAboutIndex.LabelList = *classifyList

	return nil
}

// 1.查kind和icon
func getKindAndIcon(k *[]models.DataAboutKind) error {
	return mysql.GetDataAboutKind(k)
}

// 2.查classifyDetails
func getClassifyAndDetails(c *[]models.DataAboutLabel) error {
	//得到了所有的分类
	return mysql.GetAllDataAboutClassify(c)
}

// 3.查询allEssay
func getAllEssay(data *[]models.DataAboutEssay) error {
	return mysql.GetAllEssay(data)
}

func ResponseDataAboutEssayList(essayList *[]models.DataAboutEssay) (err error) {
	if err = getAllEssay(essayList); err != nil {
		return err
	}
	// 整合classify
	var classifyList = new([]models.DataAboutLabel)
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
