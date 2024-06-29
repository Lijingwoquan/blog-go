package logic

import (
	"blog/dao/redis"
	"blog/models"
)

func SetEssayKeyword(essayDetail *models.EssayKeyword) (err error) {
	err = redis.SetEssayKeyword(essayDetail)
	return err
}

func checkKeywordExist(keyword string) (exist bool, err error) {
	return exist, err
}

func IncreaseSearchKeyword(SearchKeyword *models.Keyword) (err error) {
	// 1.先看该关键词是否存在 --> 遍历关键词数据

	// 2.再对该关键词进行处理  --> 向redis中加入keyWord
	return redis.IncreaseSearchKeyword((*SearchKeyword).Keyword)
}
