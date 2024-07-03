package redis

func SaveUserIp(ip string) (err error) {
	return SetYearMonthWeekTimesZoneForZset(ip, KeyUserIp, 1)
}
