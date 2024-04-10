package redis

import (
	"fmt"
)

func GetVisitedTimes(eid int64) (int64, error) {
	pre := getRedisKey(KeyVisitedTimes)
	vt, err := client.HIncrBy(pre, fmt.Sprintf("%d", eid), 1).Result()
	if err != nil {
		return 0, err
	}
	return vt, err
}

func InitVisitedTimes(eid int64) error {
	pre := getRedisKey(KeyVisitedTimes)
	_, err := client.HSet(pre, fmt.Sprintf("%d", eid), 0).Result()
	return err
}
