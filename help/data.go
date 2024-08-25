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

	var essayList = new([]models.DataAboutEssay)
	if err = getAllEssay(essayList); err != nil {
		return err
	}

	sortIndexData(DataAboutIndex, kindList, classifyList, essayList)

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

// 3.查询allEssay
func getAllEssay(data *[]models.DataAboutEssay) error {
	return mysql.GetAllEssay(data)
}

// 4.整合数据

func sortIndexData(DataAboutIndex *models.DataAboutIndex, k *[]models.DataAboutKind, c *[]models.DataAboutClassify, e *[]models.DataAboutEssay) {
	sortKindAndClassify(DataAboutIndex, k, c)
	sortClassifyAndEssay(DataAboutIndex, c, e)
}

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

func sortClassifyAndEssay(DataAboutIndex *models.DataAboutIndex, c *[]models.DataAboutClassify, e *[]models.DataAboutEssay) {
	var indexDataEssayList = make([]models.DataAboutEssay, 0, len(*e))

	// 计算 indexDataEssayList 的总大小，并创建具有适当容量的切片
	var classifyRouterMap = make(map[string]string, len(*c))
	for _, classify := range *c {
		classifyRouterMap[classify.Name] = classify.Router
	}

	var essayClassifyMap = make(map[string][]models.DataAboutEssay, len(*c))
	for _, essay := range *e {
		essayClassifyMap[essay.Kind] = append(essayClassifyMap[essay.Kind], essay)
	}

	for k, v := range essayClassifyMap {
		kindRoute := classifyRouterMap[k]
		for _, essay := range v {
			complexRouter := "/essay" + kindRoute + essay.Router
			indexDataEssayList = append(indexDataEssayList, models.DataAboutEssay{
				Name:          essay.Name,
				Kind:          essay.Kind,
				Router:        essay.Router,
				ComplexRouter: complexRouter,
				Introduction:  essay.Introduction,
				ID:            essay.ID,
				Keywords:      essay.Keywords,
			})
		}
	}

	if err := redis.GetEssayKeywordsForIndex(&indexDataEssayList); err != nil {
		return
	}
	DataAboutIndex.EssayList = indexDataEssayList
}
