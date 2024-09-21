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
	var classifyList = new([]models.DataAboutClassify)
	if err = getClassifyAndDetails(classifyList); err != nil {
		return err
	}

	//整合数据
	sortKindAndClassify(DataAboutIndex, kindList, classifyList)
	return nil
}

func ResponseDataAboutEssayList(essayList *[]models.DataAboutEssay) (err error) {
	if err = getAllEssay(essayList); err != nil {
		return err
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
	return mysql.GetDataAboutClassifyDetails(c)
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

	DataAboutIndex.Menu = indexDataMenu
}
