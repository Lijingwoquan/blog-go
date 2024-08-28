package logic

import (
	"blog/dao/redis"
	"blog/models"
)

func GetUserIpCount(ipSet *models.UserIpForSet) (err error) {
	return redis.GetUserIpCount(ipSet)
}

func GetSearchKeywordRank(rankKind *models.RankKindForZset) (err error) {
	//	得到年月日的keywords的zset
	return redis.GetSearchKeywordRank(rankKind)
}
