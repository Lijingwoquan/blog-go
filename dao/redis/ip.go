package redis

import (
	"blog/models"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

const (
	ExpirationIncrement = 2 * time.Second
	MaxExpirationTime   = 10 * time.Second
)

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

func IncreaseIpRequestTimes(ip string) (times int64, err error) {
	preKey := getRedisKey(KeyLimitIp)

	pipe := client.Pipeline()

	incr := pipe.HIncrBy(preKey, ip, 1)

	ttl := pipe.TTL(preKey)

	_, err = pipe.Exec()

	if err != nil {
		return 0, err
	}

	times, err = incr.Result()
	if err != nil {
		return 0, err
	}

	// 获取剩余过期时间
	remainingTime, err := ttl.Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return times, err
	}

	// 如果键不存在或已过期，设置基础过期时间
	if remainingTime < 0 {
		remainingTime = 0
	}
	// 计算新的过期时间，但不超过最大值
	newExpiration := remainingTime + ExpirationIncrement
	if newExpiration > MaxExpirationTime {
		newExpiration = MaxExpirationTime
	}

	// 设置新的过期时间
	err = client.Expire(preKey, newExpiration).Err()
	if err != nil {
		return times, err
	}

	return times, err
}

func CheckIpIfMalicious(ip string) (bool, error) {
	preKey := getRedisKey(KeyMaliciousIp)
	return client.SIsMember(preKey, ip).Result()
}

func SetIpMalicious(ip string) (err error) {
	preKey := getRedisKey(KeyMaliciousIp)
	if err := client.SAdd(preKey, ip).Err(); err != nil {
		return err
	}
	return err
}
