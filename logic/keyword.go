package logic

import (
	"blog/dao/redis"
	"blog/models"
)

func IncreaseSearchKeyword(SearchKeyword *models.KeywordParam) (err error) {
	preKey := redis.KeySearchKeyWordTimes
	// 向redis中加入keyWord
	return redis.IncreaseSearchKeyword(preKey, (*SearchKeyword).Keyword)
}

func GetSearchKeywordRank(rankKind *models.RankKindForZset) (err error) {
	//	得到年月日的keywords的zset
	return redis.GetSearchKeywordRank(rankKind)
}
