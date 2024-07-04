package redis

import "blog/models"

func SaveUserIp(ip string) (err error) {
	preKey := getRedisKey(KeyUserIp)
	return SetYearMonthWeekTimesZoneForZset(preKey, ip, 1)
}

func GetUserIpRank(rankKind *models.RankKind) (err error) {
	preKey := getRedisKey(KeyUserIp)
	return GetYearMonthWeekTimesZoneForZsetRank(rankKind, preKey)
}
