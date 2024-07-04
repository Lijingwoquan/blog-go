package redis

import "blog/models"

func SaveUserIp(ip string) {
	preKey := getRedisKey(KeyUserIp)
	if err := SetYearMonthWeekTimesZoneForSet(preKey, ip); err != nil {
		return
	}
}

func GetUserIpCount(setKind *models.SetKind) (err error) {
	preKey := getRedisKey(KeyUserIp)
	return GetYearMonthWeekTimesZoneForSetCount(setKind, preKey)
}
