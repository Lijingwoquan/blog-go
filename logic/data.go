package logic

import (
	"blog/cache"
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

func GetDataAboutIndex(data *models.DataAboutIndex) (err error) {
	// 从缓存中拿到数据
	if cache.Error != nil {
		return cache.Error
	}
	var temp = new(models.DataAboutIndex)
	if temp, err = cache.GetDataAboutIndex(); err != nil {
		return err
	}
	*data = *temp
	return nil
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
	classify := new(models.DataAboutClassify)
	classify.Name = data.Kind
	if err := mysql.GetOneDataAboutClassify(classify); err != nil {
		return err
	}
	data.KindRouter = classify.Router

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

func GetDataAboutClassifyEssayMsg(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	// 先由参数查询essay的内容
	if err := mysql.GetDataAboutClassifyEssayMsg(data, query); err != nil {
		return err
	}

	for i, essay := range *data.EssayList {
		// 1.查访问次数
		var err error
		if (*data.EssayList)[i].VisitedTimes, err = redis.GetVisitedTimes(essay.Eid); err != nil {
			return err
		}
		// 2.查关键字
		if (*data.EssayList)[i].Keywords, err = redis.GetEssayKeywordsForOne(essay.Eid); err != nil {
			return err
		}
	}

	var classifyList = new([]models.DataAboutClassify)
	if err := mysql.GetAllDataAboutClassify(classifyList); err != nil {
		return err
	}

	var classifyMap = make(map[string]string)
	for _, classify := range *classifyList {
		classifyMap[classify.Name] = classify.Router
	}
	// 再遍历essayList为它们加上kindRouter
	for i, essay := range *data.EssayList {
		(*data.EssayList)[i].KindRouter = classifyMap[essay.Kind]
	}
	return nil
}
