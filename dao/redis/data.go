package redis

import (
	"blog/dao/mysql"
	"fmt"
	"strconv"
)

func GetVisitedTimes(eid string) (int64, error) {
	var vt int64
	key := getRedisKey(KeyVisitedTimes)
	//先检查键是否存在
	var exist bool
	var err error
	if exist, err = client.HExists(key, eid).Result(); err != nil {
		return 0, err
	}
	if !exist {
		if vt, err = mysql.GetVisitedTimesFromMySQL(eid); err != nil {
			return 0, err
		}
	}

	if vt, err = client.HIncrBy(key, eid, vt+1).Result(); err != nil {
		return 0, err
	}

	//将访问的文章加入到集合中
	if _, err = client.SAdd(getRedisKey(KeyChangeVisitedTimes), eid).Result(); err != nil {
		return 0, err
	}
	return vt, nil
}

func InitVisitedTimes(eid int64) error {
	pre := getRedisKey(KeyVisitedTimes)
	_, err := client.HSet(pre, fmt.Sprintf("%d", eid), 0).Result()
	return err
}

func GetChangedVisitedTimes() (map[int64]int64, error) {
	pre := getRedisKey(KeyVisitedTimes)
	eids, err := client.SMembers(getRedisKey(KeyChangeVisitedTimes)).Result()
	if err != nil {
		return nil, err
	}
	var visitedTimesChangedMap = make(map[int64]int64)
	for _, eid := range eids {
		eida, err := strconv.ParseInt(eid, 10, 64)
		vt, err := client.HGet(pre, eid).Int64()
		if err != nil {
			return nil, err
		}
		visitedTimesChangedMap[eida] = vt
	}
	//删除这个集合
	_, err = client.Del(getRedisKey("visited:eids")).Result()
	if err != nil {
		return nil, err
	}
	return visitedTimesChangedMap, nil
}
