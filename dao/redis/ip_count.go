package redis

import "blog/models"

func SaveUserIp(ip string) {
	preKey := getRedisKey(KeyUserIp)
	if err := SetYearMonthWeekTimesZoneForSet(preKey, ip); err != nil {
		return
	}
}

func GetUserIpCount(ipSet *models.UserIpForSet) (err error) {
	preKey := getRedisKey(KeyUserIp)
	return GetYearMonthWeekTimesZoneForSet(ipSet, preKey)
}
