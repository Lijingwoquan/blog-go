package logic

import (
	"blog/cache"
	"blog/dao/redis"
	"blog/models"
	"fmt"
)

func SetEssayKeyword(essayDetail *models.EssayKeyword) (err error) {
	err = redis.SetEssayKeyword(essayDetail)
	return err
}

func checkKeywordExist(keyword string) (exist bool) {
	fmt.Println(cache.GlobalEssayKeyword)
	for _, value := range cache.GlobalEssayKeyword {
		if value == keyword {
			exist = true
			break
		}
	}
	return exist
}

func IncreaseSearchKeyword(SearchKeyword *models.Keyword) (err error) {
	// 1.先看该关键词是否存在 --> 遍历关键词数据
	exist := checkKeywordExist((*SearchKeyword).Keyword)

	fmt.Println(exist)

	// 2.再对该关键词进行处理  --> 向redis中加入keyWord
	if exist {
		return redis.IncreaseSearchKeyword((*SearchKeyword).Keyword)
	} else {
		return nil
	}
}
