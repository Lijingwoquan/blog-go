package logic

import (
	"blog/dao/redis"
	"blog/models"
)

func IncreaseSearchKeyword(SearchKeyword *models.KeywordParam) (err error) {
	// 向redis中加入keyWord
	return redis.IncreaseSearchKeyword((*SearchKeyword).Keyword)
}

func GetSearchKeywordRank(rankKind *models.KeywordRankKind) (err error) {
	//	得到年月日的keywords的zset
	return redis.GetSearchKeywordRank(rankKind)
}
