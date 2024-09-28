package logic

import (
	"blog/cache"
	"blog/dao/redis"
	"blog/models"
	"strings"
)

func GetDataByKeyword(param *models.SearchParam) (err error) {
	//判断是否需要添加set值
	if param.IfAdd {
		preKey := redis.KeySearchKeyWordTimes
		// 向redis中加入keyWord
		return redis.IncreaseSearchKeyword(preKey, (*param).Keyword)
	}
	return nil
}

func GetEssayListByKeyword(e *[]models.DataAboutEssay, k *models.SearchParam) {
	var essayList = new([]models.DataAboutEssay)
	var err error
	if essayList, err = cache.GetEssayList(); err != nil {
		return
	}
	for _, essay := range *essayList {
		// 检查 essay.keyword 数组中是否包含指定的关键字 k
		for _, keyword := range essay.Keywords {
			if strings.Contains(keyword, k.Keyword) {
				// 如果包含，将文章添加到结果列表 e 中
				*e = append(*e, essay)
				break // 找到匹配的关键字后，可以跳出内层循环
			}
		}
	}
}
