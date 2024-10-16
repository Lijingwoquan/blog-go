package help

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

// ResponseIndexData 返回首页数据
func ResponseIndexData(DataAboutIndex *models.DataAboutIndex) (err error) {
	var kindList = new([]models.DataAboutKind)
	if err = mysql.GetKindList(kindList); err != nil {
		return err
	}
	var labelList = new([]models.DataAboutLabel)
	if err = mysql.GetLabelList(labelList); err != nil {
		return err
	}

	var essayList = new([]models.DataAboutEssay)
	if err = mysql.GetRecommendEssayList(essayList); err != nil {
		return err
	}

	//整合数据
	DataAboutIndex.KindList = *kindList
	DataAboutIndex.LabelList = *labelList
	DataAboutIndex.EssayList = *essayList
	return nil
}

func ResponseDataAboutEssayList(essayList *[]models.DataAboutEssay) (err error) {
	if err = mysql.GetAllEssay(essayList); err != nil {
		return err
	}
	// 整合classify
	var labelList = new([]models.DataAboutLabel)
	if err = mysql.GetLabelList(labelList); err != nil {
		return err
	}

	return redis.GetEssayKeywordsForIndex(essayList)
}
