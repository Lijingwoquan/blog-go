package help

import (
	"blog/dao/mysql"
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

	*DataAboutIndex = *sortIndexData(kindList, classifyList)
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

// 3.整合kind数据和classify数据
// 整合逻辑
// 自下而上 先得到classify和kindName组成的map
// 遍历kindList 再向menu里面插入kind数据 此时使用上文的kindName来插入classify数据
// 核心点就在于 找到公用新key

func sortIndexData(k *[]models.DataAboutKind, c *[]models.DataAboutClassify) *models.DataAboutIndex {
	var indexData = models.DataAboutIndex{}
	var indexDataMenu = make([]models.DataAboutIndexMenu, len(*k))

	var kindAndClassifyMap = make(map[string][]models.DataAboutClassify)
	for _, classify := range *c {
		kindAndClassifyMap[classify.Kind] = append(kindAndClassifyMap[classify.Name], classify)
	}

	for i, kind := range *k {
		indexDataMenu[i].DataAboutKind = kind
		indexDataMenu[i].Classify = kindAndClassifyMap[kind.ClassifyKind]
	}

	indexData.DataAboutIndexMenu = indexDataMenu
	return &indexData
}
