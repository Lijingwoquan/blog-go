package logic

import (
	"blog/cache"
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
	"strings"
)

func GetEssayList(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	if err := mysql.GetEssayList(data, query); err != nil {
		return err
	}
	return nil
}

func GetDataByKeyword(e *[]models.DataAboutEssay, param *models.SearchParam) (err error) {
	//判断是否需要添加访问值
	if param.IfAdd {
		preKey := redis.KeySearchKeyWordTimes
		// 向redis中加入keyWord
		if err = redis.IncreaseSearchKeyword(preKey, (*param).Keyword); err != nil {
			return err
		}
	}

	var essayList = new([]models.DataAboutEssay)
	essayList = cache.GetAllEssayList()
	for _, essay := range *essayList {
		// 检查 essay.keyword 数组中是否包含指定的关键字 k
		for _, keyword := range essay.Keywords {
			if strings.Contains(keyword, param.Keyword) {
				*e = append(*e, essay)
				break
			}
		}
	}
	return nil
}
