package logic

import (
	"blog/cache"
	"blog/dao/redis"
	"blog/models"
	"strings"
)

func GetDataByKeyword(e *[]models.DataAboutEssay, param *models.SearchParam) (err error) {
	//判断是否需要添加访问值
	if param.IfAdd {
		preKey := redis.KeySearchKeyWordTimes
		// 向redis中加入keyWord
		return redis.IncreaseSearchKeyword(preKey, (*param).Keyword)
	}

	var essayList = new([]models.DataAboutEssay)
	if essayList, err = cache.GetEssayList(); err != nil {
		return
	}
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
