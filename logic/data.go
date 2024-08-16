package logic

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

// GetEssayData 得到文章数据
func GetEssayData(data *models.EssayData, id int) error {
	var err error
	//从mysql查数据
	if err = mysql.GetEssayData(data, id); err != nil {
		return err
	}
	//从redis查数据
	// 1.查访问次数
	if data.VisitedTimes, err = redis.GetVisitedTimes(data.Eid); err != nil {
		return err
	}
	// 2.查关键字
	if err = redis.GetEssayKeywordsForOne(data); err != nil {
		return err
	}
	return nil
}

func GetDataAboutClassifyEssayMsg(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	return mysql.GetDataAboutClassifyEssayMsg(data, query)
}
