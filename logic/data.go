package logic

import (
	"blog/cache"
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

func GetIndexData() (data *models.DataAboutIndex, err error) {
	// 从缓存中拿到数据
	if data, err = cache.GetIndexData(); err != nil {
		return nil, err
	}
	return data, nil
}

// GetEssayData 得到文章数据
func GetEssayData(data *models.EssayData, id int) error {
	var err error
	//从mysql查数据
	//1.拿到essay本身数据
	if err = mysql.GetEssayData(data, id); err != nil {
		return err
	}
	//2.整合classify
	classify := new(models.DataAboutLabel)
	classify.Name = data.Kind
	if err := mysql.GetOneDataAboutClassify(classify); err != nil {
		return err
	}

	//从redis查数据
	// 1.查访问次数
	if data.VisitedTimes, err = redis.GetVisitedTimes(data.Eid); err != nil {
		return err
	}
	data.VisitedTimes++
	//2.更新访问次数
	if err = redis.AddVisitedTimes(data.Eid); err != nil {
		return err
	}
	// 3.查关键字
	if data.Keywords, err = redis.GetEssayKeywordsForOne(data.Eid); err != nil {
		return err
	}
	return nil
}

func GetEssayList(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	if err := mysql.GetEssayList(data, query); err != nil {
		return err
	}
	return nil
}
