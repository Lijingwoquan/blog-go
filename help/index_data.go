package help

import (
	"blog/dao/mysql"
	"blog/models"
)

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
		return err
	}

	//整合数据
	DataAboutIndex.KindList = *kindList
	DataAboutIndex.LabelList = *labelList
	DataAboutIndex.EssayList = *essayList
	return nil
}
